package test

import (
	"bytes"
	"os"
	"testing"

	"github.com/fishtailstudio/imgo"
)

// test output image bytes data
func TestOutPutData(t *testing.T) {
	data := new(bytes.Buffer)

	// png
	img := imgo.Load("../examples/gopher.png").OutPut(data, "png")
	os.WriteFile("out.png", data.Bytes(), 0777)

	// jpg jpeg
	data = new(bytes.Buffer)
	img.OutPut(data, "jpg")
	os.WriteFile("out.jpg", data.Bytes(), 0777)

	// tiff
	data = new(bytes.Buffer)
	img.OutPut(data, "tiff")
	os.WriteFile("out.tiff", data.Bytes(), 0777)

	// bmp
	data = new(bytes.Buffer)
	img.OutPut(data, "bmp")
	os.WriteFile("out.bmp", data.Bytes(), 0777)
}
