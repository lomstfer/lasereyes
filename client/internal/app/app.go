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
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type App struct {
	startTime           float64
	time                float64
	lastUpdateLocalTime float64
	textFace            *text.GoTextFaceSource

	finishedAssetLoading bool

	netClient *network.NetworkClient

	startedClosingProcedure bool
	timeOfCloseInput        time.Time
	cleanClose              bool

	timeSyncer *network.TimeSyncer

	getInputCallback  *common.FixedCallback
	sendInputCallback *common.FixedCallback

	mousePositionScreen vec2.Vec2
	mousePositionWorld  vec2.Vec2

	cameraTopLeftPos     vec2.Vec2
	cameraTopLeftPosInit bool

	selfPlayer             *internal.SelfPlayer
	otherPlayers           map[uint]*internal.Player
	playerEyeImage         *ebiten.Image
	playerHealthBarBgImage *ebiten.Image
	playerHealthBarFgImage *ebiten.Image
	playerPupilImage       *ebiten.Image

	gridShader *ebiten.Shader

	bufferedShootInput *vec2.Vec2

	mousePositionLastSendDirection vec2.Vec2

	laserBeams     []internal.LaserBeam
	laserBeamImage *ebiten.Image
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
	app.playerHealthBarBgImage.Fill(color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	app.playerHealthBarFgImage = ebiten.NewImage(1, 1)
	app.playerHealthBarFgImage.Fill(color.NRGBA{G: 255, A: 255})

	var err error
	app.gridShader, err = ebiten.NewShader(utils.LoadBytesInFs(assetFS, "embed_assets/gridShader.kage"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app.laserBeamImage = ebiten.NewImage(1, 1)
	app.laserBeamImage.Fill(color.NRGBA{R: 255, G: 255, B: 255, A: 255})

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
	dt := localTime - a.lastUpdateLocalTime
	a.lastUpdateLocalTime = localTime

	{
		mx, my := ebiten.CursorPosition()
		a.mousePositionScreen = vec2.Vec2{X: float64(mx), Y: float64(my)}
		a.mousePositionWorld = a.mousePositionScreen.Add(a.cameraTopLeftPos)
	}

	if a.selfPlayer != nil {
		a.UpdateSelfPlayer()
		cameraTowards := vec2.NewVec2(a.selfPlayer.SmoothedPosition.X-float64(screen.Bounds().Dx())/2, a.selfPlayer.SmoothedPosition.Y-float64(screen.Bounds().Dy())/2)
		a.cameraTopLeftPos = moveCameraTowardsSmoothly(a.cameraTopLeftPos, cameraTowards, constants.CameraSpeed, dt)
	}
	for _, p := range a.otherPlayers {
		p.LerpBetweenSnapshots(a.timeSyncer.ServerTime())
	}

	laserBeamsLeft := make([]internal.LaserBeam, 0)
	for _, lb := range a.laserBeams {
		if a.time-lb.TimeInstantiated < constants.LaserBeamViewTime {
			laserBeamsLeft = append(laserBeamsLeft, lb)
		}
	}
	a.laserBeams = laserBeamsLeft

	a.draw(screen)
}

func (a *App) UpdateSelfPlayer() {
	if a.selfPlayer.Data.Dead {
		return
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && a.bufferedShootInput == nil {
		a.bufferedShootInput = &vec2.Vec2{X: a.mousePositionWorld.X, Y: a.mousePositionWorld.Y}
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
			if a.mousePositionLastSendDirection != a.mousePositionWorld {
				packetStruct := msgfromclient.UpdateFacingDirection{Dir: a.selfPlayer.Data.PupilDistDir01}
				bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromclient.MsgTypeUpdateFacingDirection), packetStruct)
				a.netClient.SendToServer(bytes, false)
			}
			a.mousePositionLastSendDirection = a.mousePositionWorld
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

	a.selfPlayer.UpdateSmoothPosition(a.getInputCallback.Accumulator / a.getInputCallback.DeltaSeconds)

	a.selfPlayer.CalculateFacingVec(a.mousePositionWorld)
}

func (a *App) loadAssets(assetFS embed.FS) {
	fontBytes, err := assetFS.ReadFile("embed_assets/Roboto-Regular.ttf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	a.textFace = utils.GetTextFace(fontBytes)

	a.finishedAssetLoading = true
}

func (a *App) syncTime() {
	for !a.timeSyncer.FinishedSync {
		request := msgfromclient.TimeRequest{TimeSent: commonutils.GetUnixTimeAsFloat()}
		bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromclient.MsgTypeTimeRequest), request)
		a.netClient.SendToServer(bytes, false)
		time.Sleep(time.Millisecond * constants.ServerTimeSyncDeltaMS)
	}
}

func (a *App) getPlayerDataFromId(id uint) *player.CommonData {
	if id == a.selfPlayer.id {
		return &a.selfPlayer.Data
	}
	p := a.otherPlayers[id]
	if p != nil {
		return &p.Data
	}
	return nil
}

func (a *App) onCloseInput() {
	a.startedClosingProcedure = true
	a.timeOfCloseInput = time.Now()
	a.netClient.StartDisconnect()
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
		pd := a.getPlayerDataFromId(msg.PlayerId)
		if pd == nil {
			break
		}
		pd.Health -= msg.Damage
		if pd.Health <= 0 {
			pd.Dead = true
		}
		a.laserBeams = append(a.laserBeams, internal.LaserBeam{TargetId: pd.Id, TimeInstantiated: a.time, OwnerId: msg.CausingDamageId})
	}
}

func (a *App) draw(screen *ebiten.Image) {
	if !a.finishedAssetLoading || !a.timeSyncer.FinishedSync {
		screen.Fill(color.NRGBA{R: 100, G: 100, B: 100, A: 255})
		f := &text.GoTextFace{
			Source: a.textFace,
			Size:   100,
		}
		width, _ := text.Measure("connecting", f, 0)
		geo := ebiten.GeoM{}
		geo.Translate(float64(screen.Bounds().Dx())/2-width/2, 100)
		text.Draw(screen, "connecting", f, &text.DrawOptions{DrawImageOptions: ebiten.DrawImageOptions{GeoM: geo}})
		return
	}

	if !a.cameraTopLeftPosInit {
		a.cameraTopLeftPos = vec2.NewVec2(-float64(screen.Bounds().Dx())/2, -float64(screen.Bounds().Dy())/2)
		a.cameraTopLeftPosInit = true
	}

	a.drawBgGrid(screen)

	cameraTranslation := getCameraTranslation(a.cameraTopLeftPos)

	for _, p := range a.otherPlayers {
		internal.DrawPlayer(p.Data, screen, a.playerEyeImage, a.playerPupilImage, cameraTranslation)
	}
	selfDataToRender := a.selfPlayer.Data
	selfDataToRender.Position = a.selfPlayer.SmoothedPosition
	internal.DrawPlayer(selfDataToRender, screen, a.playerEyeImage, a.playerPupilImage, cameraTranslation)

	for _, lb := range a.laserBeams {
		ownerData, targetData := a.getPlayerDataFromId(lb.OwnerId), a.getPlayerDataFromId(lb.TargetId)
		if ownerData | == nil {
			break
		}
		lb.Draw(screen, internal.GetPupilPos(*ownerData), , a.laserBeamImage, a.time, cameraTranslation)
	}

	// ui stuff

	playerNameTextFace := &text.GoTextFace{
		Source: a.textFace,
		Size:   20,
	}
	for _, p := range a.otherPlayers {
		internal.DrawPlayerHealthBar(p.Data, screen, a.playerHealthBarBgImage, a.playerHealthBarFgImage, cameraTranslation)
		internal.DrawPlayerName(p.Data, screen, playerNameTextFace, cameraTranslation)
	}
	internal.DrawPlayerHealthBar(selfDataToRender, screen, a.playerHealthBarBgImage, a.playerHealthBarFgImage, cameraTranslation)
	internal.DrawPlayerName(selfDataToRender, screen, playerNameTextFace, cameraTranslation)
}

func (a *App) drawBgGrid(screen *ebiten.Image) {
	opts := ebiten.DrawRectShaderOptions{}
	opts.Uniforms = make(map[string]any, 2)
	opts.Uniforms["CameraTopLeft"] = []float32{
		float32(a.cameraTopLeftPos.X),
		float32(a.cameraTopLeftPos.Y),
	}
	screen.DrawRectShader(screen.Bounds().Dx(), screen.Bounds().Dy(), a.gridShader, &opts)
}

func moveCameraTowardsSmoothly(cameraPos vec2.Vec2, towards vec2.Vec2, step float64, deltaTime float64) vec2.Vec2 {
	diff := towards.Sub(cameraPos)
	return cameraPos.Add(diff.Mul(step).Mul(deltaTime))
}

func getCameraTranslation(cameraTopLeftPosition vec2.Vec2) vec2.Vec2 {
	return cameraTopLeftPosition.Mul(-1)
}
