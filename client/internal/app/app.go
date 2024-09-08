package app

import (
	"embed"
	"fmt"
	"image/color"
	"os"
	"time"
	"wzrds/client/internal"
	"wzrds/client/internal/constants"
	"wzrds/client/internal/network"
	"wzrds/client/pkg/utils"
	"wzrds/common"
	commonconstants "wzrds/common/commonconstants"
	commonutils "wzrds/common/commonutils"
	"wzrds/common/netmsg"
	"wzrds/common/netmsg/msgfromclient"
	"wzrds/common/netmsg/msgfromserver"
	"wzrds/common/pkg/vec2"
	"wzrds/common/player"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type App struct {
	startTime           float64
	time                float64
	lastUpdateLocalTime float64
	fontFace            *font.Face

	finishedAssetLoading bool

	netClient *network.NetworkClient

	startedClosingProcedure bool
	timeOfCloseInput        time.Time
	cleanClose              bool

	timeSyncer *network.TimeSyncer

	getInputCallback  *common.FixedCallback
	sendInputCallback *common.FixedCallback

	selfPlayer             *internal.SelfPlayer
	otherPlayers           map[uint]*internal.Player
	playerEyeImage         *ebiten.Image
	playerHealthBarBgImage *ebiten.Image
	playerHealthBarFgImage *ebiten.Image
	playerPupilImage       *ebiten.Image

	backgroundImage *ebiten.Image

	bufferedShootInput *vec2.Vec2
}

func NewApp(assetFS embed.FS) *App {
	app := &App{}
	app.startTime = commonutils.GetUnixTimeAsFloat()
	app.lastUpdateLocalTime = app.startTime

	app.timeSyncer = network.NewTimeSyncer(constants.TimesToSyncClock)

	app.getInputCallback = common.NewFixedCallback(commonconstants.SimulationTickRate)

	app.sendInputCallback = common.NewFixedCallback(constants.SendInputRate)
	app.otherPlayers = make(map[uint]*internal.Player)
	app.playerEyeImage = ebiten.NewImageFromImage(*utils.LoadImageInFs(assetFS, "embed_assets/eye.png"))
	app.playerPupilImage = ebiten.NewImageFromImage(*utils.LoadImageInFs(assetFS, "embed_assets/pupil.png"))

	app.playerHealthBarBgImage = ebiten.NewImage(1, 1)
	app.playerHealthBarBgImage.Fill(color.NRGBA{255, 255, 255, 255})
	app.playerHealthBarFgImage = ebiten.NewImage(1, 1)
	app.playerHealthBarFgImage.Fill(color.NRGBA{255, 0, 0, 255})

	{
		bgImg := ebiten.NewImage(30, 30)
		dark := true
		for y := range bgImg.Bounds().Dy() {
			for x := range bgImg.Bounds().Dx() {
				if dark {
					bgImg.Set(x, y, color.NRGBA{20, 20, 20, 255})
				} else {
					bgImg.Set(x, y, color.NRGBA{60, 60, 60, 255})
				}
				dark = !dark
			}
			dark = !dark
		}
		app.backgroundImage = bgImg
	}

	app.netClient = network.NewNetworkClient()

	go app.loadAssets(assetFS)

	return app
}

func (a *App) UpdateClose() bool {
	if ebiten.IsWindowBeingClosed() && !a.startedClosingProcedure {
		a.onCloseInput()
	}

	if a.startedClosingProcedure {
		if a.cleanClose || time.Since(a.timeOfCloseInput).Seconds() > constants.WaitForCleanCloseTime {
			a.netClient.CleanUp()
			return true
		}
	}

	return false
}

func (a *App) Update(screen *ebiten.Image) {
	a.handleNetworkEvents()

	a.time = commonutils.GetUnixTimeAsFloat() - a.startTime
	localTime := commonutils.GetUnixTimeAsFloat()
	// dt := timeNow - g.lastUpdateTime
	a.lastUpdateLocalTime = localTime

	var mousePosition vec2.Vec2
	{
		mx, my := ebiten.CursorPosition()
		mousePosition = vec2.Vec2{X: float64(mx), Y: float64(my)}
	}

	if a.selfPlayer != nil {
		a.UpdateSelfPlayer(mousePosition)
	}
	for _, p := range a.otherPlayers {
		p.LerpBetweenSnapshots(a.timeSyncer.ServerTime())
	}

	a.draw(screen)
}

func (a *App) UpdateSelfPlayer(mousePosition vec2.Vec2) {
	if a.selfPlayer.Data.Dead {
		return
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && a.bufferedShootInput == nil {
		a.bufferedShootInput = &vec2.Vec2{X: mousePosition.X, Y: mousePosition.Y}
	}

	a.getInputCallback.Update(func() {
		inputVec := vec2.Vec2{}
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			inputVec.Y -= 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			inputVec.Y += 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			inputVec.X -= 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			inputVec.X += 1
		}
		if !a.selfPlayer.Data.Dead {
			a.selfPlayer.CheckMoveInput(inputVec)
		}
	})

	a.sendInputCallback.Update(func() {
		{
			packetStruct := msgfromclient.UpdateFacingDirection{Dir: a.selfPlayer.Data.PupilDistDir01}
			bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromclient.MsgTypeUpdateFacingDirection), packetStruct)
			a.netClient.SendToServer(bytes, false)
		}
		{
			if len(a.selfPlayer.InputsToSend) == 0 && a.bufferedShootInput == nil {
				return
			}
			move := msgfromclient.MoveInput{MoveInputs: a.selfPlayer.InputsToSend}
			shoot := msgfromclient.ShootInput{Time: a.timeSyncer.ServerTime()}
			if a.bufferedShootInput != nil {
				shoot.DidShoot = true
				shoot.Position = *a.bufferedShootInput
			}
			a.bufferedShootInput = nil
			packetStruct := msgfromclient.Input{Move: move, Shoot: shoot}
			bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromclient.MsgTypeInput), packetStruct)
			a.netClient.SendToServer(bytes, true)
			a.selfPlayer.OnSendInputs()
		}
	})

	a.selfPlayer.UpdateRenderPosition(a.getInputCallback.Accumulator / a.getInputCallback.DeltaSeconds)

	a.selfPlayer.CalculateFacingVec(mousePosition)
}

func (g *App) loadAssets(assetFS embed.FS) {
	fontBytes, err := assetFS.ReadFile("embed_assets/Roboto-Regular.ttf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	g.fontFace = utils.GetFontFace(fontBytes)

	g.finishedAssetLoading = true
}

func (g *App) syncTime() {
	for !g.timeSyncer.FinishedSync {
		request := msgfromclient.TimeRequest{TimeSent: commonutils.GetUnixTimeAsFloat()}
		bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromclient.MsgTypeTimeRequest), request)
		g.netClient.SendToServer(bytes, false)
		time.Sleep(time.Millisecond * constants.ServerTimeSyncDeltaMS)
	}
}

func (g *App) onCloseInput() {
	g.startedClosingProcedure = true
	g.timeOfCloseInput = time.Now()
	g.netClient.StartDisconnect()
}

func (a *App) handleNetworkEvents() {
	netMessage := a.netClient.CheckForEvents()
	switch msg := netMessage.(type) {
	case network.Connected:
		go a.syncTime()

	case network.Disconnected:
		a.cleanClose = true

	case msgfromserver.AddSelfPlayer:
		a.selfPlayer = internal.NewSelfPlayer(msg.Data)

	case msgfromserver.RemoveOtherPlayer:
		delete(a.otherPlayers, msg.Id)

	case msgfromserver.AddOtherPlayer:
		a.otherPlayers[msg.Data.Id] = &internal.Player{Data: msg.Data}

	case msgfromserver.UpdatePlayers:
		for id, snapshot := range msg.IdsToSnapshots {
			if id == a.selfPlayer.Data.Id {
				continue
			}
			a.otherPlayers[id].SnapshotsForInterp = append(a.otherPlayers[id].SnapshotsForInterp, snapshot)
		}

	case msgfromserver.UpdateSelf:
		a.selfPlayer.HandleServerUpdate(msg.LastAuthorizedInputId, msg.Snapshot)

	case msgfromserver.TimeAnswer:
		a.timeSyncer.OnTimeAnswer(commonutils.GetUnixTimeAsFloat(), msg.Request.TimeSent, msg.TimeReceived)

	case msgfromserver.PlayerTakeDamage:
		var pd *player.CommonData
		if msg.PlayerId == a.selfPlayer.Data.Id {
			pd = &a.selfPlayer.Data
		} else {
			pd = &a.otherPlayers[msg.PlayerId].Data
		}
		pd.Health -= msg.Damage
		if pd.Health <= 0 {
			pd.Dead = true
		}
	}
}

func (a *App) draw(screen *ebiten.Image) {
	if !a.finishedAssetLoading || !a.timeSyncer.FinishedSync {
		screen.Fill(color.NRGBA{100, 100, 100, 255})
		text.Draw(screen, "connecting", *a.fontFace, 0, 24, color.NRGBA{255, 255, 255, 255})
		return
	}

	{
		geo := ebiten.GeoM{}
		geo.Scale(float64(screen.Bounds().Dx())/float64(a.backgroundImage.Bounds().Dx()), float64(screen.Bounds().Dy())/float64(a.backgroundImage.Bounds().Dy()))
		screen.DrawImage(a.backgroundImage, &ebiten.DrawImageOptions{GeoM: geo})
	}

	for _, p := range a.otherPlayers {
		internal.DrawPlayer(p.Data, screen, a.playerEyeImage, a.playerPupilImage)
	}
	selfDataToRender := a.selfPlayer.Data
	selfDataToRender.Position = a.selfPlayer.RenderPosition
	internal.DrawPlayer(selfDataToRender, screen, a.playerEyeImage, a.playerPupilImage)

	// ui stuff

	for _, p := range a.otherPlayers {
		internal.DrawPlayerHealthBar(p.Data, screen, a.playerHealthBarBgImage, a.playerHealthBarFgImage)
	}
	internal.DrawPlayerHealthBar(selfDataToRender, screen, a.playerHealthBarBgImage, a.playerHealthBarFgImage)
}
