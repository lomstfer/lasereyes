package utils

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
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

func LoadBytesInFs(fs embed.FS, filePath string) []byte {
	file, err := fs.Open(filePath)
	if err != nil {
		fmt.Println("failed to open file, error:", err)
		os.Exit(1)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("failed to read file, error:", err)
		os.Exit(1)
	}

	return data
}

func GetTextFace(ttf []byte) *text.GoTextFaceSource {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(ttf))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return s
}
