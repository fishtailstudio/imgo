package main

import "github.com/fishtailstudio/imgo"

func main() {
    imgo.Load("gopher.png").
        Mosaic(5, 60, 50, 120, 100).
        Save("out.png")
}
