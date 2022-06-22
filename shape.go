package imgo

import (
	"image"
	"image/color"
)

// Pixel draws a pixel at given (x, y) coordinate with given color.
func (i *Image) Pixel(x, y int, c color.Color) *Image {
	if i.Error != nil {
		return i
	}

	i.image.Set(x, y, c)

	return i
}

// Line draws a line from (x1, y1) to (x2, y2) with given color.
// TODO: width is not working yet.
func (i *Image) Line(x1, y1, x2, y2 int, c color.Color, width ...int) *Image {
	if i.Error != nil {
		return i
	}

	w := 1
	if len(width) != 0 {
		w = width[0]
	}
	if w <= 0 {
		return i
	}

	bg := i.image

	dx := x2 - x1
	dy := y2 - y1
	if dx == 0 && dy == 0 { // Line is a point.
		bg.Set(x1, y1, c)
		i.image = bg
		return i
	} else if dx == 0 { // vertical line
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for y := y1; y <= y2; y++ {
			bg.Set(x1, y, c)
		}
	} else if dy == 0 { // horizontal line
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		for x := x1; x <= x2; x++ {
			bg.Set(x, y1, c)
		}
	} else { // diagonal line
		k := float64(dy) / float64(dx)
		if x1 > x2 {
			x1, x2 = x2, x1
			y1, y2 = y2, y1
		}
		if -1 < k && k < 1 {
			for x := x1; x <= x2; x++ {
				y := int(float64(y1) + float64(x-x1)*k)
				bg.Set(x, y, c)
			}
		} else {
			ySmaller, yBigger := y1, y2
			if y1 > y2 {
				ySmaller, yBigger = y2, y1
			}
			for y := ySmaller; y <= yBigger; y++ {
				x := int(float64(x1) + float64(y-y1)/k)
				bg.Set(x, y, c)
			}
		}
	}

	return i
}

// Circle draws a circle at given center coordinate (x, y) with given radius and color.
func (i *Image) Circle(x, y, radius int, c color.Color) *Image {
	if i.Error != nil {
		return i
	}

	if radius <= 0 || radius >= i.width || radius >= i.height {
		return i
	}

	x1 := x - radius
	y1 := y - radius
	x2 := x + radius
	y2 := y + radius

	bg := i.image

	for x3 := x1; x3 < x2; x3++ {
		for y3 := y1; y3 < y2; y3++ {
			if (x3-x)*(x3-x)+(y3-y)*(y3-y) <= radius*radius {
				bg.Set(x3, y3, c)
			}
		}
	}

	return i
}

// Rectangle draws a rectangle at given coordinate (x, y) with given width and height and color.
func (i *Image) Rectangle(x, y, width, height int, c color.Color) *Image {
	if i.Error != nil {
		return i
	}

	if width <= 0 || height <= 0 {
		return i
	}

	x1 := x
	y1 := y
	x2 := x + width
	y2 := y + height

	bg := i.image

	for x3 := x1; x3 < x2; x3++ {
		for y3 := y1; y3 < y2; y3++ {
			if image.Pt(x3, y3).In(i.image.Bounds()) {
				bg.Set(x3, y3, c)
			}
		}
	}

	return i
}

// Ellipse draws an ellipse at given center coordinate (x, y) with given width and height and color.
func (i *Image) Ellipse(x, y, width, height int, c color.Color) *Image {
	if i.Error != nil {
		return i
	}

	if width <= 0 || height <= 0 {
		return i
	}

	a := float64(width) / 2
	b := float64(height) / 2
	x1 := x - int(a)
	y1 := y - int(b)
	x2 := x + int(a)
	y2 := y + int(b)

	bg := i.image

	for x3 := x1; x3 <= x2; x3++ {
		for y3 := y1; y3 <= y2; y3++ {
			if (float64(x3)-float64(x))*(float64(x3)-float64(x))/a/a+(float64(y3)-float64(y))*(float64(y3)-float64(y))/b/b <= 1.0 {
				bg.Set(x3, y3, c)
			}
		}
	}

	return i
}
