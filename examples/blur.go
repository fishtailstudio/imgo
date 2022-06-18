package main

import "github.com/fishtailstudio/imgo"

func main() {
    imgo.Load("gopher.png").
        Blur(5).
        Save("out.png")
}
