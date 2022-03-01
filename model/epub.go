package model

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/kyicy/readimension/utility/epub"
	"github.com/mholt/archiver"
	"golang.org/x/image/draw"
	"gorm.io/gorm"
)

type Epub struct {
	gorm.Model
	Title       string `gorm:"type:varchar(255)"`
	SHA256      string `gorm:"type:varchar(255);unique;not null"`
	SizeByMB    float64
	Author      string `gorm:"type:varchar(255)"`
	HasCover    bool
	CoverFormat string
}

func (e *Epub) CoverPath() string {
	return "/covers/" + e.SHA256 + "." + e.CoverFormat
}

func (e *Epub) IsZipped() bool {
	return e.SizeByMB <= 10.0
}

func (e *Epub) StoreName() string {
	if e.IsZipped() {
		return e.SHA256 + ".epub"
	} else {
		return e.SHA256
	}
}

func Scale(src image.Image) image.Image {
	p := src.Bounds().Size()
	var rect image.Rectangle
	if p.X < 300 {
		rect = src.Bounds()
	} else {
		y := int(float64(p.Y) * (300.0 / float64(p.X)))
		rect = image.Rect(0, 0, 300, y)
	}
	dst := image.NewRGBA(rect)
	draw.ApproxBiLinear.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
	return dst
}

func NewEpub(info epub.BookInfo, workDir, fileName, storeFolder string) (*Epub, error) {
	book := info.Book()
	var coverFormat string
	if info.HasCover() {
		bytes, format, _ := info.GetCover()
		coverFormat = format

		coverPath := filepath.Join(workDir, "covers", fmt.Sprintf("%s.%s", book.Hash, format))
		file, err := os.Create(coverPath)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		// resize to width 1000 using Lanczos resampling
		// and preserve aspect ratio
		m := Scale(bytes)

		switch format {
		case "gif":
			gif.Encode(file, m, nil)
		case "jpeg":
			jpeg.Encode(file, m, nil)
		case "png":
			png.Encode(file, m)
		}
	}

	epubRecord := &Epub{
		Title:       book.Title,
		SHA256:      book.Hash,
		SizeByMB:    float64(book.FileSize) / float64(1024*1024),
		Author:      book.Author,
		HasCover:    info.HasCover(),
		CoverFormat: coverFormat,
	}

	DB.Create(epubRecord)

	if epubRecord.IsZipped() {
		return epubRecord, os.Rename(fileName, fmt.Sprintf("%s.epub", storeFolder))
	}
	return epubRecord, archiver.DefaultZip.Unarchive(fileName, storeFolder)
}
