package imgo

import "errors"

var (
    ErrSourceImageIsNil          = errors.New("source image is nil")
    ErrSourceNotSupport          = errors.New("source not support")
    ErrSourceStringIsEmpty       = errors.New("source string is empty")
    ErrSourceImageNotSupport     = errors.New("source image not support")
    ErrSaveImageFormatNotSupport = errors.New("save image format not support")
)
