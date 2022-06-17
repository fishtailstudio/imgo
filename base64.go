package imgo

import (
    "bytes"
    "encoding/base64"
    "image/png"
    "strings"
)

// ToBase64 returns the base64 encoded string of the image.
func (i Image) ToBase64() string {
    buff := bytes.NewBuffer(nil)
    err := png.Encode(buff, i.image)
    if err != nil {
        return ""
    }
    return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buff.Bytes())
}

// LoadFromBase64 loads an image from a base64 encoded string.
func LoadFromBase64(base64Str string) (i *Image) {
    i = &Image{}
    base64Str = strings.Split(base64Str, ",")[1]

    // Decode the base64 string
    decodeString, err := base64.StdEncoding.DecodeString(base64Str)
    if err != nil {
        i.addError(err)
        return
    }

    // Get the extension, mimetype and corresponding decoder function of the image.
    ext, mime, decoder, err := GetImageType(decodeString[:8])
    if err != nil {
        i.addError(err)
        return
    }

    // Decode the image.
    buff := bytes.NewBuffer(decodeString)
    img, err := decoder(buff)
    if err != nil {
        i.addError(err)
        return
    }

    return &Image{
        image:     Image2RGBA(img),
        width:     img.Bounds().Dx(),
        height:    img.Bounds().Dy(),
        extension: ext,
        mime:      mime,
    }
}
