package main

import (
	"embed"
	"log"
	"wzrds/client/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed embed_assets/Roboto-Regular.ttf
//go:embed embed_assets/dud.png
var assetFS embed.FS

func main() {
	game := game.NewGame(assetFS)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
