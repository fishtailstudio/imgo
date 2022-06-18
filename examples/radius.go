package main

import (
    "github.com/fishtailstudio/imgo"
    "image/color"
)

func main() {
    imgo.Canvas(300, 300, color.White).
        BorderRadius(20).
        Save("out.png")
}
