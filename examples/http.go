package main

import (
    "github.com/fishtailstudio/imgo"
    "net/http"
)

func main() {
    http.HandleFunc("/gopher", imgo.Load("gopher.png").HttpHandler)
    http.ListenAndServe(":8080", nil)
}
