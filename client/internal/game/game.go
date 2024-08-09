package game

import (
	"embed"
	"errors"
	"fmt"
	"image/color"
	"os"
	"time"
	"wzrds/client/internal"
	"wzrds/client/internal/constants"
	"wzrds/client/internal/network"
	"wzrds/client/pkg/utils"
	"wzrds/common"
	"wzrds/common/netmsg"
	"wzrds/common/netmsg/msgfromclient"
	commonutils "wzrds/common/utils"

	"github.com/hajimehoshi/ebiten/v2"
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

	timeSyncer *network.TimeSyncer

	getInputCallback  *common.FixedCallback
	sendInputCallback *common.FixedCallback

	selfPlayer   *internal.SelfPlayer
	otherPlayers map[uint]*internal.Player
	playerImage  *ebiten.Image
}

func NewGame(assetFS embed.FS) *Game {
	ebiten.SetWindowSize(640, 360)
	ebiten.SetWindowTitle("pocketino")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowClosingHandled(true)

	game := &Game{}
	game.startTime = time.Now()
	game.lastUpdateTime = time.Now()

	game.timeSyncer = network.NewTimeSyncer(10)

	game.getInputCallback = common.NewFixedCallback(1.0 / 60.0)

	game.sendInputCallback = common.NewFixedCallback(1.0 / 30.0)
	game.otherPlayers = make(map[uint]*internal.Player)
	game.playerImage = ebiten.NewImage(20, 20)
	game.playerImage.Fill(color.NRGBA{255, 0, 0, 255})

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

func (g *Game) syncTime(deltaSleep time.Duration) {
	for !g.timeSyncer.FinishedSync {
		request := msgfromclient.TimeRequest{TimeSent: commonutils.GetCurrentTimeAsFloat()}
		bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromclient.MsgTypeTimeRequest), request)
		g.netClient.SendToServer(bytes, true)
		time.Sleep(deltaSleep)
	}
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 500, 500
}

func (g *Game) OnCloseInput() {
	g.startedClosingProcedure = true
	g.timeOfCloseInput = time.Now()
	g.netClient.StartDisconnect()
}
