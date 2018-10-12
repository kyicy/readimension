package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/kyicy/readimension/model"
	"github.com/mholt/archiver"

	"github.com/kyicy/readimension/utility/epub"
	"github.com/labstack/echo"
)

type getBooksData struct {
	*TempalteCommon
	Books []model.Book
}

func getBooks(c echo.Context) error {
	tc := newTemplateCommon(c, "Books")
	data := &getBooksData{}
	data.TempalteCommon = tc

	var books []model.Book

	userID, err := getSessionUserID(c)
	if err != nil {
		return err
	}

	model.DB.Where("user_id = ?", userID).Preload("Epub").Find(&books)

	data.Books = books

	return c.Render(http.StatusOK, "books", data)
}

func getBooksNew(c echo.Context) error {
	tc := newTemplateCommon(c, "Books")
	data := &gtbData{}
	data.TempalteCommon = tc
	return c.Render(http.StatusOK, "books/new", data)
}

const uploadDir = "uploads"

// Request parameters
const (
	paramUUID = "qquuid" // uuid
	paramFile = "qqfile" // file name
)

// Chunked request parameters
const (
	paramPartIndex       = "qqpartindex"      // part index
	paramPartBytesOffset = "qqpartbyteoffset" // part byte offset
	paramTotalFileSize   = "qqtotalfilesize"  // total file size
	paramTotalParts      = "qqtotalparts"     // total parts
	paramFileName        = "qqfilename"       // file name for chunked requests
	paramChunkSize       = "qqchunksize"      // size of the chunks
)

type UploadResponse struct {
	Success      bool   `json:"success"`
	Error        string `json:"error,omitempty"`
	PreventRetry bool   `json:"preventRetry"`
}

func writeUploadResponse(w *echo.Response, err error) {
	uploadResponse := new(UploadResponse)
	if err != nil {
		uploadResponse.Error = err.Error()
	} else {
		uploadResponse.Success = true
	}
	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode(uploadResponse)
}

func writeHTTPResponse(w *echo.Response, httpCode int, err error) {
	w.WriteHeader(httpCode)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func postBooksNew(c echo.Context) error {
	req := c.Request()
	w := c.Response()
	uuid := req.FormValue(paramUUID)
	if len(uuid) == 0 {
		return errors.New("invalid upload request")
	}

	file, headers, err := req.FormFile(paramFile)

	if err != nil {
		writeUploadResponse(w, err)
		return nil
	}

	fileDir := fmt.Sprintf("%s/%s", uploadDir, uuid)
	if err := os.MkdirAll(fileDir, 0777); err != nil {
		writeUploadResponse(w, err)
		return nil
	}

	var filename string
	partIndex := req.FormValue(paramPartIndex)
	if len(partIndex) == 0 {
		filename = fmt.Sprintf("%s/%s", fileDir, headers.Filename)
	} else {
		filename = fmt.Sprintf("%s/%s_%05s", fileDir, uuid, partIndex)
	}

	outfile, err := os.Create(filename)
	if err != nil {
		writeUploadResponse(w, err)
		return nil
	}
	defer outfile.Close()

	if _, err := io.Copy(outfile, file); err != nil {
		writeUploadResponse(w, err)
		return nil
	}

	writeUploadResponse(w, nil)

	if len(partIndex) == 0 {
		afterUpload(c, filename)
	}
	return nil
}

func postChunksDone(c echo.Context) error {
	req := c.Request()
	w := c.Response()

	uuid := req.FormValue(paramUUID)
	filename := req.FormValue(paramFileName)
	totalFileSize, err := strconv.Atoi(req.FormValue(paramTotalFileSize))
	if err != nil {
		writeHTTPResponse(w, http.StatusInternalServerError, err)
		return nil
	}
	totalParts, err := strconv.Atoi(req.FormValue(paramTotalParts))
	if err != nil {
		writeHTTPResponse(w, http.StatusInternalServerError, err)
		return nil
	}

	finalFilename := fmt.Sprintf("%s/%s/%s", uploadDir, uuid, filename)
	f, err := os.Create(finalFilename)
	if err != nil {
		writeHTTPResponse(w, http.StatusInternalServerError, err)
		return nil
	}
	defer f.Close()

	var totalWritten int64
	for i := 0; i < totalParts; i++ {
		part := fmt.Sprintf("%[1]s/%[2]s/%[2]s_%05[3]d", uploadDir, uuid, i)
		partFile, err := os.Open(part)
		if err != nil {
			writeHTTPResponse(w, http.StatusInternalServerError, err)
			return nil
		}
		written, err := io.Copy(f, partFile)
		if err != nil {
			writeHTTPResponse(w, http.StatusInternalServerError, err)
			return nil
		}
		partFile.Close()
		totalWritten += written

		if err := os.Remove(part); err != nil {
			c.Logger().Errorf("Error: %v", err)
		}
	}

	if totalWritten != int64(totalFileSize) {
		errorMsg := fmt.Sprintf("Total file size mistmatch, expected %d bytes but actual is %d", totalFileSize, totalWritten)
		http.Error(w, errorMsg, http.StatusMethodNotAllowed)
	}

	return afterUpload(c, finalFilename)
}

func afterUpload(c echo.Context, fileName string) error {
	defer func() {
		// remove upload folder
		path := filepath.Dir(fileName)
		os.RemoveAll(path)
	}()

	info, err := epub.Load(fileName)

	// not a epub file
	if err != nil {
		return err
	}

	userRecord, err := getSessionUser(c)
	if err != nil {
		return err
	}

	book := info.Book()
	storeFolder := "books/" + book.Hash
	storeName := storeFolder + ".epub"

	var epubRecord model.Epub
	model.DB.Where("sha256 = ?", book.Hash).First(&epubRecord)

	if epubRecord.SHA256 != book.Hash {
		os.MkdirAll("covers", 0777)

		var coverFormat string
		if info.HasCover() {
			bytes, format, err := info.GetCover()
			coverFormat = format

			file, err := os.Create("covers/" + book.Hash + "." + format)
			if err != nil {
				return err
			}

			switch format {
			case "gif":
				gif.Encode(file, bytes, nil)
			case "jpeg":
				jpeg.Encode(file, bytes, nil)
			case "png":
				png.Encode(file, bytes)
			}
		}
		os.MkdirAll("books", 0777)
		os.Rename(fileName, storeName)

		epubRecord = model.Epub{
			Title:       book.Title,
			SHA256:      book.Hash,
			SizeByMB:    float64(book.FileSize) / float64(1024*1024),
			Author:      book.Author,
			HasCover:    info.HasCover(),
			CoverFormat: coverFormat,
		}

		model.DB.Create(&epubRecord)

	}

	model.DB.Model(userRecord).Association("Books").Append(
		model.Book{
			Epub: epubRecord,
		},
	)

	defer func() {
		archiver.Zip.Open(storeName, storeFolder)
		os.Remove(storeName)
	}()

	return nil
}
