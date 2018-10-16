package route

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/kyicy/readimension/model"
	"github.com/labstack/echo"
)

func getExplorerRoot(c echo.Context) error {
	userID, err := getSessionUserID(c)
	if err != nil {
		return err
	}
	var user model.User
	model.DB.Where("id = ?", userID).Find(&user)

	return c.Redirect(http.StatusFound, fmt.Sprintf("/u/explorer/%v", user.ListID))
}

type getBooksData struct {
	*TempalteCommon
	List model.List
}

func getExplorer(c echo.Context) error {
	id := c.Param("id")

	tc := newTemplateCommon(c, "Library Explorer")
	data := &getBooksData{}
	data.TempalteCommon = tc
	data.Active = "/u/explorer"

	var list model.List
	model.DB.Where("id = ?", id).Preload("Epubs").Preload("Children").Find(&list)
	data.List = list

	return c.Render(http.StatusOK, "explorer", data)
}

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

	newList := model.List{
		Name:     _list.Name,
		ParentID: uint(id),
	}

	model.DB.Create(&newList)
	model.DB.Model(listRecord).Association("Children").Append(newList)

	r := make(map[string]string)
	r["id"] = fmt.Sprintf("%v", newList.ID)
	r["name"] = newList.Name

	return c.JSON(http.StatusOK, r)
}
