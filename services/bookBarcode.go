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

func (agent DBAgent) HasBook(isbn string, idOptional ...int) bool {
	var (
		err   error
		value int
	)
	if len(idOptional) > 0 {
		id := idOptional[0]
		row := agent.DB.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT * from book WHERE id='%v' AND isbn='%v');", id, isbn))
		err = row.Scan(&value)
	} else {
		row := agent.DB.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT * from book WHERE isbn='%v');", isbn))
		err = row.Scan(&value)
	}
	if err != nil || value == 0 {
		return false
	}
	return true
}

func EscapeForSQL(sql string) string {
	dest := make([]byte, 0, 2*len(sql))
	var escape byte
	for i := 0; i < len(sql); i++ {
		c := sql[i]

		escape = 0

		switch c {
		case 0: /* Must be escaped for 'mysql' */
			escape = '0'
			break
		case '\n': /* Must be escaped for logs */
			escape = 'n'
			break
		case '\r':
			escape = 'r'
			break
		case '\\':
			escape = '\\'
			break
		case '\'':
			escape = '\''
			break
		case '"': /* Better safe than sorry */
			escape = '"'
			break
		case '\032': //十进制26,八进制32,十六进制1a, /* This gives problems on Win32 */
			escape = 'Z'
		}

		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}

	return string(dest)
}

func (agent *DBAgent) AddBookBarcode(id int, isbn string) *StatusResult {
	// 首先检查 id以及isbn是否在数据库中
	if !agent.HasBook(isbn, id) {
		return &StatusResult{
			Msg:    "数据库中不存在该书籍",
			Status: BookBarcodeFailed,
		}
	}

	var value int
	codingMsg := fmt.Sprintf("%v-%v", isbn, id)
	savePath := filepath.Join(mediaPath, fmt.Sprintf("%v.png", codingMsg))

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
	var pngFile *os.File

	pngFile, err = os.Create(savePath)
	if err != nil {
		return &StatusResult{
			Msg:    "文件创建失败:" + err.Error(),
			Status: BookBarcodeFailed,
		}
	}
	defer pngFile.Close()

	err = png.Encode(pngFile, img)
	if err != nil {
		return &StatusResult{
			Msg:    "png编码失败:" + err.Error(),
			Status: BookBarcodeFailed,
		}
	}

	row := agent.DB.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT * from book WHERE isbn='%v');", isbn))
	if temperr := row.Scan(&value); temperr == nil && value != 0 {
		if temp, tempPath := agent.GetBookBarcodePath(id, isbn); temp.Status == BookBarcodeOK && tempPath == savePath {
			return &StatusResult{
				Msg:    "数据库中已经存在该书籍对应条形码",
				Status: BookBarcodeOK,
			}
		} else {
			result, sqlerr := agent.DB.Exec(fmt.Sprintf(`UPDATE book_barcode
			SET barcode_path = '%v'
			WHERE id='%v' AND isbn='%v';`,
				EscapeForSQL(savePath), id, isbn))
			if sqlerr != nil {
				return &StatusResult{
					Msg:    "SQL存储失败: " + sqlerr.Error(),
					Status: BookBarcodeFailed,
				}
			}
			if noOfRow, temperr := result.RowsAffected(); temperr != nil || noOfRow <= 0 {
				return &StatusResult{
					Msg:    "SQL存储失败: " + temperr.Error(),
					Status: BookBarcodeFailed,
				}
			}
			return &StatusResult{
				Msg:    "编码成功",
				Status: BookBarcodeOK,
			}
		}
	}
	result, sqlerr := agent.DB.Exec(fmt.Sprintf(`INSERT INTO book_barcode(id,isbn,barcode_path) 
			VALUES ('%v','%v','%v')`,
		id, isbn, EscapeForSQL(savePath)))
	if sqlerr != nil {
		return &StatusResult{
			Msg:    "SQL存储失败: " + sqlerr.Error(),
			Status: BookBarcodeFailed,
		}
	}
	if noOfRow, temperr := result.RowsAffected(); temperr != nil || noOfRow <= 0 {
		return &StatusResult{
			Msg:    "SQL存储失败: " + temperr.Error(),
			Status: BookBarcodeFailed,
		}
	}
	return &StatusResult{
		Msg:    "编码成功",
		Status: BookBarcodeOK,
	}
}
