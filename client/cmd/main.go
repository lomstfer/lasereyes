package main

import (
	"embed"
	_ "embed"
	"log"
	"wzrds/client/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed embed_assets/Roboto-Regular.ttf
var assetFS embed.FS

func main() {
	game := game.NewGame(assetFS)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
