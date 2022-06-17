package main

import (
	"image/color"

	"github.com/fishtailstudio/imgo"
)

func main() {
	imgo.Canvas(300, 300, color.White).
		RadiusBorder(20).
		Save("out.png")
}
