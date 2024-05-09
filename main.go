package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/brandonmpark/raytracer/internal/display"
	"github.com/brandonmpark/raytracer/internal/reader"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: raytracer /path/to/scene")
		return
	}
	scene := reader.ReadFile(os.Args[1])
	screen := display.Draw(scene)

	img := image.NewNRGBA(image.Rect(0, 0, scene.Width, scene.Height))
	for i := range screen {
		for j := range screen[i] {
			r := uint8(screen[i][j].X() * 255)
			g := uint8(screen[i][j].Y() * 255)
			b := uint8(screen[i][j].Z() * 255)
			img.Set(i, scene.Height-j, color.NRGBA{R: r, G: g, B: b, A: 255})
		}
	}

	file, _ := os.Create(scene.OutputFile)
	defer file.Close()
	png.Encode(file, img)
}
