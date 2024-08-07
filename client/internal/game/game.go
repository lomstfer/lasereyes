package game

import (
	"embed"
	"errors"
	"fmt"
	"image/color"
	"os"
	"time"
	"wzrds/client/internal/constants"
	"wzrds/client/internal/network"
	"wzrds/client/pkg/utils"
	"wzrds/common"
	"wzrds/common/netmsg"
	"wzrds/common/netmsg/msgfromclient"
	"wzrds/common/netmsg/msgfromserver"
	"wzrds/common/pkg/vec2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Game struct {
	startTime      time.Time
	time           float64 // seconds
	lastUpdateTime time.Time
	fontFace       *font.Face

	finishedAssetLoading bool

	netClient *network.NetworkClient

	startedClosingProcedure bool
	timeOfCloseInput        time.Time
	cleanClose              bool

	getInputCallback        *common.FixedCallback
	sendInputCallback       *common.FixedCallback
	accumulatedPlayerInputs []common.PlayerInput
}

func NewGame(assetFS embed.FS) *Game {
	ebiten.SetWindowSize(640, 360)
	ebiten.SetWindowTitle("pocketino")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowClosingHandled(true)

	game := &Game{}
	game.startTime = time.Now()
	game.lastUpdateTime = time.Now()

	game.getInputCallback = common.NewFixedCallback(1.0 / 60.0)

	game.sendInputCallback = common.NewFixedCallback(1.0 / 30.0)
	game.accumulatedPlayerInputs = make([]common.PlayerInput, 0)

	game.netClient = network.NewNetworkClient()

	go game.loadAssets(assetFS)

	return game
}

func (g *Game) loadAssets(assetFS embed.FS) {
	fontBytes, err := assetFS.ReadFile("embed_assets/Roboto-Regular.ttf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	g.fontFace = utils.GetFontFace(fontBytes)

	g.finishedAssetLoading = true
}

// called 60 times per second
func (g *Game) Update() error {
	if ebiten.IsWindowBeingClosed() && !g.startedClosingProcedure {
		g.OnCloseInput()
	}

	if g.startedClosingProcedure {
		if g.cleanClose || time.Since(g.timeOfCloseInput).Seconds() > constants.WaitForCleanCloseTime {
			g.netClient.CleanUp()
			return errors.New("window closed")
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.finishedAssetLoading {
		screen.Fill(color.NRGBA{255, 0, 0, 255})
		return
	}

	netMessage := g.netClient.CheckForEvents()
	switch /* msg :=  */ netMessage.(type) {
	case msgfromserver.DisconnectSelf:
		g.cleanClose = true
	}

	g.time = time.Since(g.startTime).Seconds()
	timeNow := time.Now()
	// dt := timeNow.Sub(g.lastUpdateTime).Seconds()
	g.lastUpdateTime = timeNow

	g.getInputCallback.Update(func() {
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
		fmt.Println(inputVec)
		input := common.PlayerInput{Up: inputVec.Y == -1, Down: inputVec.Y == 1, Left: inputVec.X == -1, Right: inputVec.X == 1}
		g.accumulatedPlayerInputs = append(g.accumulatedPlayerInputs, input)
	})

	g.sendInputCallback.Update(func() {
		packetStruct := msgfromclient.MoveInput{Input: g.accumulatedPlayerInputs}
		g.accumulatedPlayerInputs = make([]common.PlayerInput, 0)
		bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromclient.MsgTypeMoveInput), packetStruct)
		g.netClient.SendToServer(bytes, true)
	})

	text.Draw(screen, "hello", *g.fontFace, 0, 24, color.NRGBA{255, 255, 255, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 500, 500
}

func (g *Game) OnCloseInput() {
	g.startedClosingProcedure = true
	g.timeOfCloseInput = time.Now()
	g.netClient.StartDisconnect()
}
