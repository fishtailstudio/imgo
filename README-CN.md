# ImGo
[English](README.md) | 简体中文

![Github stars](https://img.shields.io/github/stars/fishtailstudio/imgo?style=social)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/fishtailstudio/imgo)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/fishtailstudio/imgo)
[![GoDoc](https://godoc.org/github.com/fishtailstudio/imgo?status.svg)](https://pkg.go.dev/github.com/fishtailstudio/imgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/fishtailstudio/imgo)](https://goreportcard.com/report/github.com/fishtailstudio/imgo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 介绍

`Image Golang` => `Img Go` => `ImGo` `[ˈɪmɡəʊ]`

`ImGo` 是一个开源的 Golang 图片处理和操作的库。它为创建、编辑和合成图像提供了一种更简单、更具表现力的方法。

## 安装

```bash
go get -u github.com/fishtailstudio/imgo
```

## 文档

[English Documentation](https://imgo.gitbook.io/en/) | [简体中文文档](https://imgo.gitbook.io/cn/)

## 如何使用

```go
package main

import "github.com/fishtailstudio/imgo"

func main() {
    imgo.Load("background.png").
        Resize(250, 350).
        Insert("gopher.png", 50, 50).
        Save("out.png")
}
```

## 维护者

[@fishtailstudio](https://github.com/fishtailstudio)

## 如何贡献

非常欢迎大家 [提交 Issue](https://github.com/fishtailstudio/imgo/issues/new) 或 [提交 Pull Request](https://github.com/fishtailstudio/imgo/pulls)。

## 点个 Star ! ⭐

如果你喜欢或正在使用这个项目，麻烦你点个 Star，谢谢！

## 使用许可

[MIT](LICENSE) © fishtailstudio