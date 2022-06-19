package imgo

import (
    "fmt"
    "golang.org/x/image/bmp"
    "golang.org/x/image/tiff"
    "golang.org/x/image/webp"
    "image"
    "image/color"
    "image/draw"
    "image/jpeg"
    "image/png"
    "io"
)

// GetImageType returns the extension, mimetype and corresponding decoder function of the image.
// It judges the image by its first few bytes called magic number.
func GetImageType(bytes []byte) (ext string, mimetype string, decoder func(r io.Reader) (image.Image, error), err error) {
    if len(bytes) < 2 {
        err = ErrSourceImageNotSupport
        return
    }

    if bytes[0] == 0xFF && bytes[1] == 0xD8 {
        ext = "jpg"
        mimetype = "image/jpeg"
        decoder = jpeg.Decode
    }

    if len(bytes) >= 4 && bytes[0] == 0x89 && bytes[1] == 0x50 && bytes[2] == 0x4E && bytes[3] == 0x47 {
        ext = "png"
        mimetype = "image/png"
        decoder = png.Decode
    }

    if bytes[0] == 0x42 && bytes[1] == 0x4D {
        ext = "bmp"
        mimetype = "image/x-ms-bmp"
        decoder = bmp.Decode
    }

    if (bytes[0] == 0x49 && bytes[1] == 0x49) || (bytes[0] == 0x4D && bytes[1] == 0x4D) {
        ext = "tiff"
        mimetype = "image/tiff"
        decoder = tiff.Decode
    }

    if bytes[0] == 0x52 && bytes[1] == 0x49 {
        ext = "webp"
        mimetype = "image/webp"
        decoder = webp.Decode
    }

    /*if bytes[0] == 0x47 && bytes[1] == 0x49 && bytes[2] == 0x46 && bytes[3] == 0x38 {
          ext = "gif"
          mimetype = "image/gif"
      }

      if bytes[0] == 0x00 && bytes[1] == 0x00 && (bytes[2] == 0x01 || bytes[2] == 0x02) && bytes[3] == 0x00 {
          ext = "ico"
          mimetype = "image/x-icon"
      }*/

    if ext == "" {
        err = ErrSourceImageNotSupport
    }

    return
}

// Color2Hex converts a color.Color to its hex string representation.
func Color2Hex(c color.Color) string {
    r, g, b, _ := c.RGBA()
    return fmt.Sprintf("#%02X%02X%02X", uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

// Image2RGBA converts an image to RGBA.
func Image2RGBA(img image.Image) *image.RGBA {
    rgba := image.NewRGBA(img.Bounds())
    draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Over)
    return rgba
}
