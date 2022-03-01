package epub

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"image"
	_ "image/gif"  // supports gif
	_ "image/jpeg" // supports jpeg
	_ "image/png"  // supports png
	"io"
	"os"
	"path/filepath"

	"github.com/beevik/etree"
	"github.com/pkg/errors"
	"golang.org/x/tools/godoc/vfs/zipfs"
)

type Book struct {
	Hash     string
	FilePath string
	FileSize int64

	HasCover bool
	Title    string
	Author   string
}

type epub struct {
	hascover  bool
	book      *Book
	coverpath *string
}

func (e *epub) Book() *Book {
	return e.book
}

func (e *epub) HasCover() bool {
	return e.coverpath != nil
}

type BookInfo interface {
	Book() *Book
	HasCover() bool
	GetCover() (image.Image, string, error)
}

func (e *epub) GetCover() (i image.Image, format string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("panic while decoding cover image")
		}
	}()

	zr, err := zip.OpenReader(e.book.FilePath)
	if err != nil {
		return nil, "", errors.Wrap(err, "error opening epub as zip")
	}
	defer zr.Close()

	zfs := zipfs.New(zr, "epub")

	cr, err := zfs.Open(*e.coverpath)

	if err != nil {
		return nil, "", errors.Wrapf(err, "could not open cover '%s'", *e.coverpath)
	}
	defer cr.Close()

	i, format, err = image.Decode(cr)
	if err != nil {
		return nil, "", errors.Wrap(err, "error decoding image")
	}

	return i, format, nil
}

func Load(filename string) (BookInfo, error) {
	e := &epub{book: &Book{}, hascover: false}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, errors.Wrapf(err, "could not stat book")
	}
	e.book.FilePath = filename
	e.book.FileSize = fi.Size()

	s := sha256.New()
	i, err := io.Copy(s, f)
	if err == nil && i != fi.Size() {
		err = errors.New("could not read whole file")
	}
	if err != nil {
		f.Close()
		return nil, errors.Wrap(err, "could not hash book")
	}
	e.book.Hash = fmt.Sprintf("%x", s.Sum(nil))

	f.Close()

	zr, err := zip.OpenReader(filename)
	if err != nil {
		return nil, errors.Wrap(err, "error opening epub as zip")
	}
	defer zr.Close()

	zfs := zipfs.New(zr, "epub")

	rsk, err := zfs.Open("/META-INF/container.xml")
	if err != nil {
		return nil, errors.Wrap(err, "error reading container.xml")
	}
	defer rsk.Close()

	container := etree.NewDocument()
	_, err = container.ReadFrom(rsk)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing container.xml")
	}

	rootfile := ""
	for _, e := range container.FindElements("//rootfiles/rootfile[@full-path]") {
		rootfile = e.SelectAttrValue("full-path", "")
	}

	if rootfile == "" {
		return nil, errors.Wrap(err, "could not find rootfile in container.xml")
	}

	opfdir := filepath.Dir(rootfile)

	rrsk, err := zfs.Open("/" + rootfile)
	if err != nil {
		return nil, errors.Wrap(err, "error reading rootfile")
	}
	defer rrsk.Close()

	opf := etree.NewDocument()
	_, err = opf.ReadFrom(rrsk)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing rootfile")
	}

	e.book.Title = filepath.Base(e.book.FilePath)
	for _, el := range opf.FindElements("//title") {
		e.book.Title = el.Text()
		break
	}
	for _, el := range opf.FindElements("//creator") {
		e.book.Author = el.Text()
		break
	}

	for _, el := range opf.FindElements("//meta[@name='cover']") {
		coverid := el.SelectAttrValue("content", "")
		if coverid != "" {
			for _, f := range opf.FindElements("//[@id='" + coverid + "']") {
				coverPath := f.SelectAttrValue("href", "")
				if coverPath != "" {
					coverPath = "/" + opfdir + "/" + coverPath
					e.coverpath = &coverPath
				}
			}
			break
		}
	}

	return e, nil
}
