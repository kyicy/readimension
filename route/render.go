package route

import (
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/markbates/pkger"
)

// Template Struct
type Template struct {
	templates *template.Template
}

// Render object
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func getRender() *Template {
	var tmpl *template.Template
	pkger.Walk("/route/template", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
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

		f, err := pkger.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		bs, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		tt.Parse(string(bs))
		return nil
	})

	t := &Template{
		templates: tmpl,
	}

	return t
}
