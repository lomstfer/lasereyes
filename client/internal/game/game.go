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
}

func NewGame(assetFS embed.FS) *Game {
	ebiten.SetWindowSize(640, 360)
	ebiten.SetWindowTitle("pocketino")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowClosingHandled(true)

	game := &Game{}
	game.startTime = time.Now()
	game.lastUpdateTime = time.Now()

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
	case common.ServerDisconnectedClient:
		g.cleanClose = true
	}

	g.time = time.Since(g.startTime).Seconds()
	timeNow := time.Now()
	// dt := timeNow.Sub(game.lastUpdateTime).Seconds()
	g.lastUpdateTime = timeNow

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
