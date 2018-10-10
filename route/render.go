package route

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/gobuffalo/packr"
	"github.com/labstack/echo"
)

// Template Struct
type Template struct {
	templates *template.Template
}

// Render object
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var box packr.Box

func init() {
	box = packr.NewBox("./template")
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

		tt.Parse(box.String(path))
		return nil
	})

	t := &Template{
		templates: tmpl,
	}

	return t
}
