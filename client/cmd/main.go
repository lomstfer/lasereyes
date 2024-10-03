package main

import (
	"embed"
	"errors"
	"log"
	"wzrds/client/internal/app"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed embed_assets/Roboto-Regular.ttf
//go:embed embed_assets/eye.png
//go:embed embed_assets/pupil.png
//go:embed embed_assets/gridShader.kage
var assetFS embed.FS

func main() {
	ebiten.SetWindowSize(640, 360)
	ebiten.SetWindowTitle("")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowClosingHandled(true)
	ebiten.SetTPS(ebiten.SyncWithFPS)

	game := EbitenGame{app: app.NewApp(assetFS)}
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

type EbitenGame struct {
	app *app.App
}

func (eg *EbitenGame) Update() error {
	updateClose := eg.app.UpdateClose()
	if updateClose {
		return errors.New("window closed")
	}
	return nil
}

func (eg *EbitenGame) Draw(screen *ebiten.Image) {
	eg.app.Update(screen)
}

func (eg *EbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 800
}
