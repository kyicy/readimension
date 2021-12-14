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
	"github.com/kyicy/readimension/utility/config"
	"github.com/kyicy/readimension/utility/epub"
	"github.com/labstack/echo/v4"
	"github.com/mholt/archiver"
	"github.com/nfnt/resize"
)

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
	conf := config.Get()

	fileDir := filepath.Join(conf.WorkDir, uploadDir, uuid)
	if err := os.MkdirAll(fileDir, 0777); err != nil {
		writeUploadResponse(w, err)
		return nil
	}

	var filename string
	partIndex := req.FormValue(paramPartIndex)
	if len(partIndex) == 0 {
		filename = filepath.Join(fileDir, headers.Filename)
	} else {
		filename = filepath.Join(fileDir, fmt.Sprintf("%s_%05s", uuid, partIndex))
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
	conf := config.Get()

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

	finalFilename := filepath.Join(conf.WorkDir, uploadDir, uuid, filename)
	f, err := os.Create(finalFilename)
	if err != nil {
		writeHTTPResponse(w, http.StatusInternalServerError, err)
		return nil
	}

	var totalWritten int64
	for i := 0; i < totalParts; i++ {
		part := filepath.Join(conf.WorkDir, uploadDir, uuid, fmt.Sprintf("%s_%05d", uuid, i))
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

	f.Close()
	afterUpload(c, finalFilename)

	writeHTTPResponse(w, http.StatusOK, nil)
	return nil
}

func afterUpload(c echo.Context, fileName string) error {
	listID := c.Param("list_id")

	info, err := epub.Load(fileName)
	// not a epub file
	if err != nil {
		return err
	}

	book := info.Book()
	conf := config.Get()
	storeFolder := filepath.Join(conf.WorkDir, "books", book.Hash)

	var epubRecord model.Epub
	model.DB.Where("sha256 = ?", book.Hash).First(&epubRecord)

	if epubRecord.SHA256 != book.Hash {
		var coverFormat string
		if info.HasCover() {
			bytes, format, _ := info.GetCover()
			coverFormat = format

			coverPath := filepath.Join(conf.WorkDir, "covers", fmt.Sprintf("%s.%s", book.Hash, format))
			file, err := os.Create(coverPath)
			if err != nil {
				return err
			}

			defer file.Close()

			// resize to width 1000 using Lanczos resampling
			// and preserve aspect ratio
			m := resize.Resize(300, 0, bytes, resize.Lanczos3)

			switch format {
			case "gif":
				gif.Encode(file, m, nil)
			case "jpeg":
				jpeg.Encode(file, m, nil)
			case "png":
				png.Encode(file, m)
			}
		}

		epubRecord = model.Epub{
			Title:       book.Title,
			SHA256:      book.Hash,
			SizeByMB:    float64(book.FileSize) / float64(1024*1024),
			Author:      book.Author,
			HasCover:    info.HasCover(),
			CoverFormat: coverFormat,
		}

		model.DB.Create(&epubRecord)

		var err error
		if epubRecord.IsZipped() {
			err = os.Rename(fileName, fmt.Sprintf("%s.epub", storeFolder))
		} else {
			err = archiver.DefaultZip.Unarchive(fileName, storeFolder)
		}

		if err != nil {
			c.Logger().Error(err)
		}
	}

	var list model.List

	model.DB.Where("id = ?", listID).Find(&list)

	model.DB.Model(&list).Association("Epubs").Append(&epubRecord)

	userIDStr, _ := getSessionUserID(c)

	userID, _ := strconv.Atoi(userIDStr)
	ule := model.UserListEpub{
		UserID: uint(userID),
		ListID: list.ID,
		EpubID: epubRecord.ID,
	}

	model.DB.Create(&ule)

	path := filepath.Dir(fileName)
	if err := os.RemoveAll(path); err != nil {
		c.Logger().Error(err)
	}

	return nil
}
