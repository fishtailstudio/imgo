package imgo

import (
    "github.com/BurntSushi/graphics-go/graphics"
    "image"
    "image/color"
    "sync"
)

// GaussianBlur returns a blurred image.
// ksize is Gaussian kernel size
// sigma is Gaussian kernel standard deviation.
func (i *Image) GaussianBlur(ksize int, sigma float64) *Image {
    if i.Error != nil {
        return i
    }

    dst := image.NewRGBA(i.image.Bounds())
    err := graphics.Blur(dst, i.image, &graphics.BlurOptions{
        StdDev: sigma,
        Size:   ksize,
    })

    if err != nil {
        i.addError(err)
        return i
    }

    i.image = dst
    i.width = dst.Bounds().Dx()
    i.height = dst.Bounds().Dy()

    return i
}

// normalizeKernel normalizes a kernel.
func (i Image) normalizeKernel(kernel [][]float64) {
    var sum float64
    for _, row := range kernel {
        for _, value := range row {
            sum += value
        }
    }
    for i, row := range kernel {
        for j, value := range row {
            kernel[i][j] = value / sum
        }
    }
}

// Blur returns a blurred image.
// ksize is filter kernel size, it must be a odd number.
func (i *Image) Blur(ksize int) *Image {
    if i.Error != nil {
        return i
    }

    if ksize < 2 {
        return i
    }

    if ksize%2 == 0 {
        ksize++
    }

    kernel := make([][]float64, ksize)
    for p := range kernel {
        row := make([]float64, ksize)
        for q := 0; q < ksize; q++ {
            row[q] = 1
        }
        kernel[p] = row
    }

    i.imageFilterFast(kernel)

    return i
}

// imageFilterFast returns a filtered image with Goroutine.
func (i *Image) imageFilterFast(kernel [][]float64) *Image {
    i.normalizeKernel(kernel)
    dst := image.NewRGBA(i.image.Bounds())
    kernelSize := len(kernel)
    radius := (kernelSize - 1) / 2
    wg := sync.WaitGroup{}
    convolution := func(x, y int, wg *sync.WaitGroup) {
        defer wg.Done()
        var sumR, sumG, sumB, sumA uint16
        for p := -radius; p <= radius; p++ {
            for q := -radius; q <= radius; q++ {
                trueX := x + q
                trueY := y + p
                if image.Pt(trueX, trueY).In(i.Bounds()) {
                    thisR, thisG, thisB, thisA := i.image.At(trueX, trueY).RGBA()
                    sumR += uint16(float64(thisR) * kernel[q+radius][p+radius])
                    sumG += uint16(float64(thisG) * kernel[q+radius][p+radius])
                    sumB += uint16(float64(thisB) * kernel[q+radius][p+radius])
                    sumA += uint16(float64(thisA) * kernel[q+radius][p+radius])
                }
            }
        }
        dst.SetRGBA64(x, y, color.RGBA64{R: sumR, G: sumG, B: sumB, A: sumA})
    }
    wg.Add(i.width * i.height)
    for y := 0; y < i.height; y++ {
        for x := 0; x < i.width; x++ {
            go convolution(x, y, &wg)
        }
    }
    wg.Wait()
    i.image = dst
    return i
}

// imageFilter returns a filtered image.
func (i *Image) imageFilter(kernel [][]float64) *Image {
    i.normalizeKernel(kernel)
    dst := image.NewNRGBA64(i.image.Bounds())
    kernelSize := len(kernel)
    radius := (kernelSize - 1) / 2
    for y := 0; y < i.height; y++ {
        for x := 0; x < i.width; x++ {
            var sumR, sumG, sumB, sumA uint16
            for p := -radius; p <= radius; p++ {
                for q := -radius; q <= radius; q++ {
                    trueX := x + q
                    trueY := y + p
                    if image.Pt(trueX, trueY).In(i.Bounds()) {
                        thisR, thisG, thisB, thisA := i.image.At(trueX, trueY).RGBA()
                        sumR += uint16(float64(thisR) * kernel[q+radius][p+radius])
                        sumG += uint16(float64(thisG) * kernel[q+radius][p+radius])
                        sumB += uint16(float64(thisB) * kernel[q+radius][p+radius])
                        sumA += uint16(float64(thisA) * kernel[q+radius][p+radius])
                    }
                }
            }
            dst.SetRGBA64(x, y, color.RGBA64{R: sumR, G: sumG, B: sumB, A: sumA})
        }
    }
    i.image = Image2RGBA(dst)
    return i
}
