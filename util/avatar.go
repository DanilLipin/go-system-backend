package util

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/fogleman/gg"
)

var colors = []color.Color{
	color.RGBA{255, 0, 0, 255},    // Red
	color.RGBA{0, 255, 0, 255},    // Green
	color.RGBA{0, 0, 255, 255},    // Blue
	color.RGBA{255, 255, 0, 255},  // Yellow
	color.RGBA{255, 165, 0, 255},  // Orange
	color.RGBA{128, 0, 128, 255},  // Purple
	color.RGBA{0, 255, 255, 255},  // Cyan
	color.RGBA{255, 20, 147, 255}, // Deep Pink
}

func randomColor() color.Color {
	return colors[rand.Intn(len(colors))]
}

func CreateAvatar(firstName, lastName, filePath string) error {
	const S = 256

	dc := gg.NewContext(S, S)

	dc.SetColor(randomColor())
	dc.Clear()

	initials := fmt.Sprintf("%s%s", string(firstName[0:2]), string(lastName[0:2]))

	if err := dc.LoadFontFace("./fonts/Oswald-Medium.ttf", 96); err != nil {
		return err
	}

	dc.SetColor(color.White)
	dc.DrawStringAnchored(initials, S/2, S/2, 0.5, 0.5)
	dc.SavePNG(filePath)
	return nil
}
