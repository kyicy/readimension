package middleware

import (
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
)

var m *minify.M

func init() {
	m = minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("application/javascript", js.Minify)
	m.AddFunc("application/json", json.Minify)
}

type pipeWriter interface {
	io.Writer
	Close() error
}

type minifyResonseWriter struct {
	Writer pipeWriter
	http.ResponseWriter
	echo.Context
}

// Minify bla bla
func Minify(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		res := c.Response()
		rw := res.Writer
		mrw := &minifyResonseWriter{
			ResponseWriter: rw,
			Context:        c,
		}

		defer func() {
			mrw.Close()
		}()

		res.Writer = mrw
		return next(c)
	}
}

func (w *minifyResonseWriter) Write(b []byte) (int, error) {
	if w.Writer == nil {
		fullContentType := w.Context.Response().Header().Get(echo.HeaderContentType)
		basicContentType := strings.Split(fullContentType, ";")[0]
		w.Writer = m.Writer(basicContentType, w.ResponseWriter)
	}
	return w.Writer.Write(b)
}

func (w *minifyResonseWriter) Close() error {
	if w.Writer != nil {
		return w.Writer.Close()
	}
	return nil
}
