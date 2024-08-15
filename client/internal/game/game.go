package game

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"time"
	"wzrds/client/internal"
	"wzrds/client/internal/constants"
	"wzrds/client/internal/network"
	"wzrds/client/pkg/utils"
	"wzrds/common"
	commonutils "wzrds/common/commonutils"
	commonconstants "wzrds/common/constants"
	"wzrds/common/netmsg"
	"wzrds/common/netmsg/msgfromclient"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type Game struct {
	startTime      float64
	time           float64
	lastUpdateTime float64
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

	li int32
}

func NewGame(assetFS embed.FS) *Game {
	ebiten.SetWindowSize(640, 360)
	ebiten.SetWindowTitle("pocketino")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowClosingHandled(true)

	game := &Game{}
	game.startTime = commonutils.GetCurrentTimeAsFloat()
	game.lastUpdateTime = game.startTime

	game.timeSyncer = network.NewTimeSyncer(constants.TimesToSyncClock)

	game.getInputCallback = common.NewFixedCallback(commonconstants.SimulationTickRate)

	game.sendInputCallback = common.NewFixedCallback(constants.SendInputRate)
	game.otherPlayers = make(map[uint]*internal.Player)
	assetFS.ReadFile("embed_assets/dud.png")
	game.playerImage = ebiten.NewImageFromImage(*utils.LoadImageInFs(assetFS, "embed_assets/dud.png"))

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

func (g *Game) syncTime() {
	for !g.timeSyncer.FinishedSync {
		request := msgfromclient.TimeRequest{TimeSent: commonutils.GetCurrentTimeAsFloat()}
		bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromclient.MsgTypeTimeRequest), request)
		g.netClient.SendToServer(bytes, false)
		time.Sleep(time.Millisecond * constants.ServerTimeSyncDeltaMS)
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
