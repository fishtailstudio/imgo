package main

import (
	"image/color"

	"github.com/fishtailstudio/imgo"
	"golang.org/x/image/colornames"
)

func main() {
	imgo.Canvas(500, 500, color.Black).
		Pixel(10, 10, colornames.Blueviolet).
		Line(20, 10, 50, 450, colornames.Gold).
		Circle(200, 100, 50, colornames.Aqua).
		Rectangle(150, 200, 100, 150, colornames.Darkblue).
		Ellipse(400, 200, 150, 50, colornames.Tomato).
		Save("out.png")
}
