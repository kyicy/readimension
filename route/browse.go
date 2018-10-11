package route

import (
	"net/http"

	"github.com/labstack/echo"
)

type gtbData struct {
	*TempalteCommon
}

func getStream(c echo.Context) error {
	tc := newTemplateCommon(c, "Stream")
	data := &gtbData{}
	data.TempalteCommon = tc
	return c.Render(http.StatusOK, "topBooks", data)
}

func getBooks(c echo.Context) error {
	tc := newTemplateCommon(c, "Books")
	data := &gtbData{}
	data.TempalteCommon = tc
	return c.Render(http.StatusOK, "topBooks", data)
}

func getBooksNew(c echo.Context) error {
	tc := newTemplateCommon(c, "Books")
	data := &gtbData{}
	data.TempalteCommon = tc
	return c.Render(http.StatusOK, "books/new", data)
}

func postBooksNew(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func getLists(c echo.Context) error {
	tc := newTemplateCommon(c, "Books")
	data := &gtbData{}
	data.TempalteCommon = tc
	return c.Render(http.StatusOK, "topBooks", data)
}

func getListsNew(c echo.Context) error {
	tc := newTemplateCommon(c, "Books")
	data := &gtbData{}
	data.TempalteCommon = tc
	return c.Render(http.StatusOK, "topBooks", data)
}
