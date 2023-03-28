package imgo

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/BurntSushi/graphics-go/graphics"
	cliColor "github.com/fatih/color"
	"github.com/golang/freetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/tiff"
)

type Image struct {
	Error       error
	image       *image.RGBA // the image
	width       int         // image width
	height      int         // image height
	extension   string      // image extension
	mimetype    string      // image mimetype
	filesize    int64       // image filesize
	isGrayscale bool        // is grayscale image
}

// ToImage returns the instance of image.Image of the image.
func (i Image) ToImage() image.Image {
	return i.image
}

// String returns the image as a string.
func (i Image) String() string {
	return fmt.Sprintf("Extension: %v\nMimetype: %v\nWidth: %v\nHeight: %v\n", i.extension, i.mimetype, i.width, i.height)
}

// addError adds an error to imgo.
// if OnlyReason is true, only the error message is returned.
func (i *Image) addError(err error, OnlyReason ...bool) {
	log.SetPrefix("[IMGO] ")

	var onlyReason bool
	if len(OnlyReason) > 0 {
		onlyReason = OnlyReason[0]
	}

	yellow := cliColor.New(cliColor.FgYellow).SprintFunc()
	magenta := cliColor.New(cliColor.FgMagenta).SprintFunc()

	_, file, line, ok := runtime.Caller(1)

	if i.Error == nil {
		if ok && !onlyReason {
			i.Error = fmt.Errorf("%v %v", yellow(file, ":", line), magenta("Error: ", err.Error()))
		} else {
			i.Error = err
		}
	} else if err != nil {
		if ok && !onlyReason {
			i.Error = fmt.Errorf("%v\n%v %v", i.Error, yellow(file, ":", line), magenta("Error: ", err.Error()))
		} else {
			i.Error = fmt.Errorf("%v\n%v", i.Error, err)
		}
	}
}

// Extension returns the extension of the image.
func (i Image) Extension() string {
	return i.extension
}

// Mimetype returns the mimetype of the image.
func (i Image) Mimetype() string {
	return i.mimetype
}

// Height returns the height of the image.
func (i Image) Height() int {
	return i.height
}

// Width returns the width of the image.
func (i Image) Width() int {
	return i.width
}

// Filesize returns the filesize of the image, if instance is initiated from an actual file.
func (i Image) Filesize() int64 {
	return i.filesize
}

// Bounds returns the bounds of the image.
func (i Image) Bounds() image.Rectangle {
	return i.image.Bounds()
}

// Insert inserts source into the image at given (x, y) coordinate.
// source can be a file path, a URL, a base64 encoded string, an *os.File, an image.Image,
// a byte slice or an *Image.
func (i *Image) Insert(source interface{}, x, y int) *Image {
	// the image to insert
	insert := &Image{}
	switch source.(type) {
	case *Image:
		insert = source.(*Image)
	default:
		insert = Load(source)
	}

	// check errors
	if insert.Error != nil {
		i.addError(insert.Error, true)
		return i
	}
	if i.Error != nil {
		return i
	}

	// check x and y is within the image
	if x > i.width || y > i.height {
		return i
	}

	// insert the image
	draw.Draw(i.image, i.Bounds(), insert.image, insert.image.Bounds().Min.Sub(image.Pt(x, y)), draw.Over)

	return i
}

// Output image data as bytes.
// Only png, jpeg, jpg, tiff and bmp image file type are supported.
// data is output save's image data.
// imageType is the out image file type.
// quality is the quality of the image, between 1 and 100, default is 100, and is only used for jpeg images.
func (i *Image) OutPut(data *bytes.Buffer, imageType string, quality ...int) *Image {
	var err error
	if i.Error != nil {
		log.Println(i.Error)
		return i
	}

	// check imageType
	if !(imageType == "png" || imageType == "jpg" || imageType == "jpeg" || imageType == "tiff" || imageType == "bmp") {
		i.addError(ErrSaveImageFormatNotSupport)
		log.Println(i.Error)
		return i
	}

	// get the image
	var img image.Image
	if i.isGrayscale {
		// grayscale image
		gray := image.NewGray(i.image.Bounds())
		for x := 0; x < i.width; x++ {
			for y := 0; y < i.height; y++ {
				rgbColor := i.image.At(x, y)
				grayColor := gray.ColorModel().Convert(rgbColor)
				gray.Set(x, y, grayColor)
			}
		}
		img = gray
	} else {
		// RGBA image
		img = i.image
	}

	// save image to file
	if imageType == "png" {
		err = png.Encode(data, img)
	} else if imageType == "jpg" || imageType == "jpeg" {
		if len(quality) > 0 && quality[0] > 0 && quality[0] < 100 {
			err = jpeg.Encode(data, img, &jpeg.Options{Quality: quality[0]})
		} else {
			err = jpeg.Encode(data, img, &jpeg.Options{Quality: 100})
		}
	} else if imageType == "tiff" {
		err = tiff.Encode(data, img, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
	} else if imageType == "bmp" {
		err = bmp.Encode(data, img)
	}

	if err != nil {
		i.addError(err)
		log.Println(i.Error)
		return i
	}
	return i
}

// Save saves the image to the specified path.
// Only png, jpeg, jpg, tiff and bmp extensions are supported.
// path is the path the image will be saved to.
// quality is the quality of the image, between 1 and 100, default is 100, and is only used for jpeg images.
func (i *Image) Save(path string, quality ...int) *Image {
	if i.Error != nil {
		log.Println(i.Error)
		return i
	}

	// check extension
	pathSplit := strings.Split(path, ".")
	extension := pathSplit[len(pathSplit)-1]
	if !(extension == "png" || extension == "jpg" || extension == "jpeg" || extension == "tiff" || extension == "bmp") {
		i.addError(ErrSaveImageFormatNotSupport)
		log.Println(i.Error)
		return i
	}

	// create file
	file, err := os.Create(path)
	if err != nil {
		i.addError(err)
		log.Println(i.Error)
		return i
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	// get the image
	var img image.Image
	if i.isGrayscale { // grayscale image
		gray := image.NewGray(i.image.Bounds())
		for x := 0; x < i.width; x++ {
			for y := 0; y < i.height; y++ {
				rgbColor := i.image.At(x, y)
				grayColor := gray.ColorModel().Convert(rgbColor)
				gray.Set(x, y, grayColor)
			}
		}
		img = gray
	} else { // RGBA image
		img = i.image
	}

	// save image to file
	if extension == "png" {
		err = png.Encode(file, img)
	} else if extension == "jpg" || extension == "jpeg" {
		if len(quality) > 0 && quality[0] > 0 && quality[0] < 100 {
			err = jpeg.Encode(file, img, &jpeg.Options{Quality: quality[0]})
		} else {
			err = jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
		}
	} else if extension == "tiff" {
		err = tiff.Encode(file, img, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
	} else if extension == "bmp" {
		err = bmp.Encode(file, img)
	}

	if err != nil {
		i.addError(err)
		log.Println(i.Error)
		return i
	}
	return i
}

// Resize resizes the image to the specified width and height.
func (i *Image) Resize(width, height int) *Image {
	if i.Error != nil {
		return i
	}

	if width == i.width || height == i.height || (width == 0 && height == 0) {
		return i
	}

	resized := resize.Resize(uint(width), uint(height), i.image, resize.Lanczos3)
	i.image = Image2RGBA(resized)
	i.width = resized.Bounds().Dx()
	i.height = resized.Bounds().Dy()

	return i
}

// Crop Cut out a rectangular part of the current image with given width and height.
func (i *Image) Crop(x, y, width, height int) *Image {
	if i.Error != nil {
		return i
	}

	if width == i.width || height == i.height || width == 0 || height == 0 || x > i.width || y > i.height {
		return i
	}

	x1 := x + width
	y1 := y + height
	clipped := i.image.SubImage(image.Rect(x, y, x1, y1))
	i.image = Image2RGBA(clipped)
	i.width = width
	i.height = height

	return i
}

// Rotate rotates the image clockwise by the specified angle.
func (i *Image) Rotate(angle int) *Image {
	if i.Error != nil {
		return i
	}
	angle %= 360
	if angle == 0 {
		return i
	}

	// angle to radian
	radian := float64(angle) * math.Pi / 180.0
	cos := math.Cos(radian)
	sin := math.Sin(radian)

	w := float64(i.width)
	h := float64(i.height)

	// calculate the new image size
	W := int(math.Max(math.Abs(w*cos-h*sin), math.Abs(w*cos+h*sin)))
	H := int(math.Max(math.Abs(w*sin-h*cos), math.Abs(w*sin+h*cos)))

	// rotate the image
	dst := image.NewRGBA(image.Rect(0, 0, W, H))
	err := graphics.Rotate(dst, i.image, &graphics.RotateOptions{Angle: radian})
	if err != nil {
		i.addError(err)
		return i
	}

	i.image = dst
	i.width = W
	i.height = H
	return i
}

// Grayscale converts the image to grayscale.
func (i *Image) Grayscale() *Image {
	if i.Error != nil {
		return i
	}

	for x := 0; x < i.width; x++ {
		for y := 0; y < i.height; y++ {
			rgbColor := i.image.At(x, y)
			grayColor := color.GrayModel.Convert(rgbColor)
			i.image.Set(x, y, grayColor)
		}
	}

	i.isGrayscale = true

	return i
}

// Flip mirror the image vertically or horizontally.
func (i *Image) Flip(flipType FlipType) *Image {
	if i.Error != nil {
		return i
	}

	if flipType == Horizontal {
		i.flipHorizontally()
	} else if flipType == Vertical {
		i.flipVertically()
	}

	return i
}

// flipHorizontally flips the image horizontally.
func (i *Image) flipHorizontally() {
	for x := 0; x < i.width/2; x++ {
		for y := 0; y < i.height; y++ {
			pixel := i.image.At(x, y)
			i.image.Set(x, y, i.image.At(i.width-x-1, y))
			i.image.Set(i.width-x-1, y, pixel)
		}
	}
}

// flipVertically flips the image vertically.
func (i *Image) flipVertically() {
	for y := 0; y < i.height/2; y++ {
		for x := 0; x < i.width; x++ {
			pixel := i.image.At(x, y)
			i.image.Set(x, y, i.image.At(x, i.height-y-1))
			i.image.Set(x, i.height-y-1, pixel)
		}
	}
}

// Text write a text string to the image at given (x, y) coordinate.
func (i *Image) Text(label string, x, y int, fontPath string, fontColor color.Color, fontSize float64, dpi float64) *Image {
	if i.Error != nil {
		return i
	}

	// Load font
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		i.addError(err)
		return i
	}
	myFont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		i.addError(err)
		return i
	}

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(myFont)
	c.SetFontSize(fontSize)
	c.SetClip(i.Bounds())
	c.SetDst(i.image)
	uni := image.NewUniform(fontColor)
	c.SetSrc(uni)
	c.SetHinting(font.HintingNone)

	// Draw text
	pt := freetype.Pt(x, y+int(c.PointToFixed(fontSize)>>6))
	if _, err := c.DrawString(label, pt); err != nil {
		i.addError(err)
		return i
	}

	return i
}

// Thumbnail returns a thumbnail of the image with given width and height.
func (i *Image) Thumbnail(width, height int) *Image {
	if i.Error != nil {
		return i
	}

	if width >= i.width || height >= i.height || width == 0 || height == 0 {
		return i
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	err := graphics.Thumbnail(dst, i.image)
	if err != nil {
		i.addError(err)
		return i
	}

	i.image = dst
	i.width = width
	i.height = height
	return i
}

// HttpHandler responds the image as an HTTP handler.
func (i Image) HttpHandler(w http.ResponseWriter, r *http.Request) {
	if i.Error != nil {
		log.Println(i.Error)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/png")
	err := png.Encode(w, i.image)

	if err != nil {
		log.Println(err)
	}
}
