package route

import (
	"net/http"

	"github.com/kyicy/readimension/model"
	"github.com/labstack/echo"
)

type list struct {
	Name string `json:"name" validate:"required,max=100"`
}

func postListChildNew(c echo.Context) error {
	_list := new(list)
	if err := c.Bind(_list); err != nil {
		return err
	}

	if err := validate.Struct(_list); err != nil {
		return err
	}

	var listRecord model.List
	id := c.Param("id")
	model.DB.Where("id = ?", id).Find(&listRecord)
	model.DB.Model(listRecord).Association("Children").Append(model.List{
		Name: _list.Name,
	})

	return c.String(http.StatusOK, "")
}
