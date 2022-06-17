package imgo

import (
    "image"
    "image/color"
    "image/draw"
)

// BorderRadius draws rounded corners on the image with given radius.
func (i *Image) BorderRadius(radius float64) *Image {
    if radius > float64(i.width/2) || radius > float64(i.height/2) {
        return i
    }

    c := Radius{p: image.Point{X: i.width, Y: i.height}, r: int(radius)}
    dst := image.NewRGBA(i.image.Bounds())
    draw.DrawMask(dst, dst.Bounds(), i.image, image.Point{}, &c, image.Point{}, draw.Over)

    i.image = dst
    return i
}

type Radius struct {
    p image.Point // the right-bottom Point of the image
    r int         // radius
}

func (c *Radius) ColorModel() color.Model {
    return color.AlphaModel
}

func (c *Radius) Bounds() image.Rectangle {
    return image.Rect(0, 0, c.p.X, c.p.Y)
}

func (c *Radius) At(x, y int) color.Color {
    var xx, yy, rr float64
    var inArea bool

    // left up
    if x <= c.r && y <= c.r {
        xx, yy, rr = float64(c.r-x)+0.5, float64(y-c.r)+0.5, float64(c.r)
        inArea = true
    }

    // right up
    if x >= (c.p.X-c.r) && y <= c.r {
        xx, yy, rr = float64(x-(c.p.X-c.r))+0.5, float64(y-c.r)+0.5, float64(c.r)
        inArea = true
    }

    // left bottom
    if x <= c.r && y >= (c.p.Y-c.r) {
        xx, yy, rr = float64(c.r-x)+0.5, float64(y-(c.p.Y-c.r))+0.5, float64(c.r)
        inArea = true
    }

    // right bottom
    if x >= (c.p.X-c.r) && y >= (c.p.Y-c.r) {
        xx, yy, rr = float64(x-(c.p.X-c.r))+0.5, float64(y-(c.p.Y-c.r))+0.5, float64(c.r)
        inArea = true
    }

    if inArea && xx*xx+yy*yy >= rr*rr {
        return color.Alpha{}
    }
    return color.Alpha{A: 255}
}
