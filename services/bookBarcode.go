package services

import (
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

func subtitleBarcode(bc barcode.Barcode) image.Image {
	fontFace := basicfont.Face7x13
	fontColor := color.RGBA{A: 255}
	margin := 5 // Space between barcode and text

	// Get the bounds of the string
	bounds, _ := font.BoundString(fontFace, bc.Content())

	widthTxt := int((bounds.Max.X - bounds.Min.X) / 64)
	heightTxt := int((bounds.Max.Y - bounds.Min.Y) / 64)

	// calc width and height
	width := widthTxt
	if bc.Bounds().Dx() > width {
		width = bc.Bounds().Dx()
	}
	height := heightTxt + bc.Bounds().Dy() + margin

	// create result img
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// draw the barcode
	draw.Draw(img, image.Rect(0, 0, bc.Bounds().Dx(), bc.Bounds().Dy()), bc, bc.Bounds().Min, draw.Over)

	// TextPt
	offsetY := bc.Bounds().Dy() + margin - int(bounds.Min.Y/64)
	offsetX := (width - widthTxt) / 2

	point := fixed.Point26_6{
		X: fixed.Int26_6(offsetX * 64),
		Y: fixed.Int26_6(offsetY * 64),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(fontColor),
		Face: fontFace,
		Dot:  point,
	}
	d.DrawString(bc.Content())
	return img
}

func (agent *DBAgent) AddBookBarcode(id int, isbn string) *StatusResult {
	codingMsg := fmt.Sprintf("%v-%v", isbn, id)
	var code barcode.Barcode
	var err error
	code, err = code128.Encode(codingMsg)
	if err != nil {
		return &StatusResult{
			Msg:    "编码失败",
			Status: BookBarcodeFailed,
		}
	}
	code, err = barcode.Scale(code, 500, 100)
	img := subtitleBarcode(code)
	var pngfile *os.File
	pngfile, err = os.Create(filepath.Join(mediaPath, fmt.Sprintf("%v.png", codingMsg)))
	if err != nil {
		return &StatusResult{
			Msg:    "文件创建失败:" + err.Error(),
			Status: BookBarcodeFailed,
		}
	}
	defer pngfile.Close()

	err = png.Encode(pngfile, img)
	if err != nil {
		return &StatusResult{
			Msg:    "png编码失败:" + err.Error(),
			Status: BookBarcodeFailed,
		}
	}
	return &StatusResult{
		Msg:    "编码成功",
		Status: BookBarcodeOK,
	}
}
