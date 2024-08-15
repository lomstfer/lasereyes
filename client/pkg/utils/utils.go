package utils

import (
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func LoadImage(filePath string) *image.Image {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("failed to open file, error:", err)
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("failed to decode image, error:", err)
		os.Exit(1)
	}

	return &img
}

func LoadImageInFs(fs embed.FS, filePath string) *image.Image {
	file, err := fs.Open(filePath)
	if err != nil {
		fmt.Println("failed to open file, error:", err)
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("failed to decode image, error:", err)
		os.Exit(1)
	}

	return &img
}

func GetFontFace(ttf []byte) *font.Face {
	var mplusNormalFont font.Face

	tt, err := opentype.Parse(ttf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &mplusNormalFont
}
