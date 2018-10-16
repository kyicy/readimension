package route

import (
	"fmt"
	"net/http"
	"strconv"

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
	id, _ := strconv.Atoi(c.Param("id"))
	model.DB.Where("id = ?", id).Find(&listRecord)

	userIDStr, _ := getSessionUserID(c)
	userID, _ := strconv.Atoi(userIDStr)

	newList := model.List{
		Name:     _list.Name,
		ParentID: uint(id),
		User:     uint(userID),
	}

	model.DB.Create(&newList)
	model.DB.Model(listRecord).Association("Children").Append(newList)

	r := make(map[string]string)
	r["id"] = fmt.Sprintf("%v", newList.ID)
	r["name"] = newList.Name

	return c.JSON(http.StatusOK, r)
}
