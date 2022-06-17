package imgo

import (
    "image"
    "image/color"
    "image/draw"
    "math"
)

// calculateMeanAverageColorWithRect returns the mean average color of the image in given rectangle.
func (i Image) calculateMeanAverageColorWithRect(rect image.Rectangle, useSquaredAverage bool) (red, green, blue uint8) {
    var redSum float64
    var greenSum float64
    var blueSum float64

    for x := rect.Min.X; x <= rect.Max.X; x++ {
        for y := rect.Min.Y; y <= rect.Max.Y; y++ {
            pixel := i.image.At(x, y)
            col := color.RGBAModel.Convert(pixel).(color.RGBA)

            if useSquaredAverage {
                redSum += float64(col.R) * float64(col.R)
                greenSum += float64(col.G) * float64(col.G)
                blueSum += float64(col.B) * float64(col.B)
            } else {
                redSum += float64(col.R)
                greenSum += float64(col.G)
                blueSum += float64(col.B)
            }
        }
    }

    rectArea := float64((rect.Dx() + 1) * (rect.Dy() + 1))

    if useSquaredAverage {
        red = uint8(math.Round(math.Sqrt(redSum / rectArea)))
        green = uint8(math.Round(math.Sqrt(greenSum / rectArea)))
        blue = uint8(math.Round(math.Sqrt(blueSum / rectArea)))
    } else {
        red = uint8(math.Round(redSum / rectArea))
        green = uint8(math.Round(greenSum / rectArea))
        blue = uint8(math.Round(blueSum / rectArea))
    }

    return
}

// pixelate apply pixelation filter to the image in given rectangle.
func (i *Image) pixelate(size int, rectangle image.Rectangle) *Image {
    bg := image.NewRGBA(i.image.Bounds())
    draw.Draw(bg, bg.Bounds(), i.image, i.image.Bounds().Min, draw.Over)

    for x := rectangle.Min.X; x < rectangle.Max.X; x += size {
        for y := rectangle.Min.Y; y < rectangle.Max.Y; y += size {
            rect := image.Rect(x, y, x+size, y+size)

            if rect.Max.X > i.width {
                rect.Max.X = i.width
            }
            if rect.Max.Y > i.height {
                rect.Max.Y = i.height
            }

            r, g, b := i.calculateMeanAverageColorWithRect(rect, true)
            col := color.RGBA{R: r, G: g, B: b, A: 255}

            for x2 := rect.Min.X; x2 < rect.Max.X; x2++ {
                for y2 := rect.Min.Y; y2 < rect.Max.Y; y2++ {
                    bg.Set(x2, y2, col)
                }
            }
        }
    }

    i.image = bg
    return i
}

// Pixelate apply pixelation filter to the image.
// size is the size of the pixel.
func (i *Image) Pixelate(size int) *Image {
    if i.Error != nil {
        return i
    }

    if size <= 1 {
        return i
    }

    if i.width > i.height {
        if size > i.width {
            size = i.width
        }
    } else {
        if size > i.height {
            size = i.height
        }
    }

    return i.pixelate(size, i.image.Bounds())
}

// Mosaic apply mosaic filter to the image in given rectangle that
// (x1, y1) and (x2, y2) are the top-left and bottom-right coordinates of the rectangle.
// size is the size of the pixel.
func (i *Image) Mosaic(size, x1, y1, x2, y2 int) *Image {
    if i.Error != nil {
        return i
    }

    if x1 < 0 {
        x1 = 0
    }
    if y1 < 0 {
        y1 = 0
    }
    if x2 > i.width {
        x2 = i.width
    }
    if y2 > i.height {
        y2 = i.height
    }

    return i.pixelate(size, image.Rect(x1, y1, x2, y2))
}
