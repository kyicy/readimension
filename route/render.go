package route

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/gobuffalo/packr/v2"
	"github.com/labstack/echo/v4"
)

// Template Struct
type Template struct {
	templates *template.Template
}

// Render object
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var box *packr.Box

func init() {
	box = packr.New("templateBox", "./template")
}

func getRender() *Template {
	var tmpl *template.Template
	box.Walk(func(path string, f packr.File) error {
		name := filepath.Base(path)
		var tt *template.Template

		if tmpl == nil {
			tmpl = template.New(name)
		}

		if name == tmpl.Name() {
			tt = tmpl
		} else {
			tt = tmpl.New(name)
		}

		s, _ := box.FindString(path)
		tt.Parse(s)
		return nil
	})

	t := &Template{
		templates: tmpl,
	}

	return t
}
