package route

import (
	"net/http"

	"github.com/labstack/echo"
)

type gtbData struct {
	*TempalteCommon
}

func getTopBooks(c echo.Context) error {
	tc := newTemplateCommon(c, "Top Books", "TopBooks")
	data := &gtbData{}
	data.TempalteCommon = tc

	return c.Render(http.StatusOK, "topBooks", data)
}

type gdData struct {
	*TempalteCommon
}

func getDiscover(c echo.Context) error {
	tc := newTemplateCommon(c, "Discover", "Discover")
	data := &gtbData{}
	data.TempalteCommon = tc
	return c.Render(http.StatusOK, "topBooks", data)
}

type gcData struct {
	*TempalteCommon
}

func getCategories(c echo.Context) error {
	tc := newTemplateCommon(c, "Categories", "Categories")
	data := &gtbData{}
	data.TempalteCommon = tc
	return c.Render(http.StatusOK, "topBooks", data)
}
