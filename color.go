package imgo

import (
    "image/color"
    "log"
)

// PickColor returns the color of the pixel at (x, y).
func (i Image) PickColor(x, y int) (res color.RGBA) {
    if i.Error != nil {
        return
    }

    if x < 0 || x > i.width || y < 0 || y > i.height {
        return
    }
    pixel := i.image.At(x, y)
    return color.RGBAModel.Convert(pixel).(color.RGBA)
}

// MainColor returns the main color in the image
func (i *Image) MainColor() (res color.RGBA) {
    if i.Error != nil {
        log.Println(i.Error)
        return
    }

    red, green, blue := i.calculateMeanAverageColorWithRect(i.image.Bounds(), true)

    return color.RGBA{
        R: red,
        G: green,
        B: blue,
        A: 0,
    }
}
