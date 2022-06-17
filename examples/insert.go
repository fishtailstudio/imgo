package main

import (
	"image/color"

	"github.com/fishtailstudio/imgo"
)

func main() {
	imgo.Canvas(500, 500, color.White).
		Insert("gopher.png", 100, 100).
		Save("out.png")
}
