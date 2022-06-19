package imgo

import (
    "bytes"
    "image"
    "image/color"
    "io/ioutil"
    "net/http"
    "os"
)

type ImageManager struct {
}

// Load an image from source.
// source can be a file path, a URL, a base64 encoded string, an *os.File, an image.Image or a byte slice.
func Load(source interface{}) *Image {
    switch source.(type) {
    case string:
        return loadFromString(source.(string))
    case *os.File:
        return LoadFromFile(source.(*os.File))
    case image.Image:
        return LoadFromImage(source.(image.Image))
    case []byte:
        return loadFromString(string(source.([]byte)))
    case *Image:
        return LoadFromImgo(source.(*Image))
    default:
        i := &Image{}
        i.addError(ErrSourceNotSupport)
        return i
    }
}

// loadFromString loads an image when the source is a string.
func loadFromString(source string) (i *Image) {
    i = &Image{}

    if len(source) == 0 {
        i.addError(ErrSourceStringIsEmpty)
        return
    }

    if len(source) > 4 && source[:4] == "http" {
        return LoadFromUrl(source)
    } else if len(source) > 10 && source[:10] == "data:image" {
        return LoadFromBase64(source)
    } else {
        return LoadFromPath(source)
    }
}

// LoadFromUrl loads an image when the source is an url.
func LoadFromUrl(url string) (i *Image) {
    i = &Image{}

    // Get the image response from the url.
    resp, err := http.Get(url)
    if err != nil {
        i.addError(err)
        return
    }
    defer resp.Body.Close()

    // Read the image data from the response.
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    // Get the extension, mimetype and corresponding decoder function of the image.
    ext, mime, decoder, err := GetImageType(bodyBytes[:8])
    if err != nil {
        i.addError(err)
        return
    }

    // Decode the image.
    file := bytes.NewReader(bodyBytes)
    img, err := decoder(file)
    if err != nil {
        i.addError(ErrSourceNotSupport)
        return
    }

    return &Image{
        image:     Image2RGBA(img),
        width:     img.Bounds().Dx(),
        height:    img.Bounds().Dy(),
        extension: ext,
        mimetype:  mime,
    }
}

// LoadFromPath loads an image from a path.
func LoadFromPath(path string) (i *Image) {
    i = &Image{}

    file, err := os.Open(path)
    if err != nil {
        i.addError(err)
        return
    }
    defer file.Close()

    return LoadFromFile(file)
}

// LoadFromFile loads an image from a file.
func LoadFromFile(file *os.File) (i *Image) {
    i = &Image{}

    // Read the first 8 bytes of the image.
    buf := make([]byte, 8)
    _, err := file.Read(buf)
    if err != nil {
        i.addError(err)
        return
    }

    // After reading the first 8 bytes, we seek back to the beginning of the file.
    _, err = file.Seek(0, 0)
    if err != nil {
        i.addError(err)
        return
    }

    // Get the extension, mimetype and corresponding decoder function of the image.
    ext, mime, decoder, err := GetImageType(buf)
    if err != nil {
        i.addError(err)
        return
    }

    // Decode the image.
    img, err := decoder(file)
    if err != nil {
        i.addError(err)
        return
    }

    // Set the image properties.
    stat, _ := file.Stat()

    return &Image{
        image:     Image2RGBA(img),
        width:     img.Bounds().Dx(),
        height:    img.Bounds().Dy(),
        extension: ext,
        mimetype:  mime,
        filesize:  stat.Size(),
    }
}

// LoadFromImage loads an image from an instance of image.Image.
func LoadFromImage(img image.Image) (i *Image) {
    i = &Image{}

    if img == nil {
        i.addError(ErrSourceImageIsNil)
        return
    }

    var formatName string
    switch img.(type) {
    case *image.NRGBA: // png
        formatName = "png"
    case *image.RGBA: // bmp, tiff
        formatName = "png"
    case *image.YCbCr: // jpeg, webp
        formatName = "jpg"
    default:
        i.addError(ErrSourceImageNotSupport)
        return
    }

    return &Image{
        image:     Image2RGBA(img),
        width:     img.Bounds().Dx(),
        height:    img.Bounds().Dy(),
        extension: formatName,
        mimetype:  "image/" + formatName,
    }
}

// LoadFromImgo loads an image from an instance of Image.
func LoadFromImgo(i *Image) *Image {
    return i
}

// Canvas create a new empty image.
func Canvas(width, height int, fillColor ...color.Color) *Image {
    var c color.Color
    if len(fillColor) == 0 {
        c = color.Transparent
    } else {
        c = fillColor[0]
    }

    img := image.NewRGBA(image.Rect(0, 0, width, height))
    for x := 0; x < img.Bounds().Dx(); x++ {
        for y := 0; y < img.Bounds().Dy(); y++ {
            img.Set(x, y, c)
        }
    }

    return &Image{
        image:     img,
        width:     img.Bounds().Dx(),
        height:    img.Bounds().Dy(),
        extension: "png",
        mimetype:  "image/png",
    }
}
