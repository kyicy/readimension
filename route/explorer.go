package route

import (
	"fmt"
	"net/http"

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
	id := c.Param("list_id")

	tc := newTemplateCommon(c, "Library Explorer")
	data := &getBooksData{}
	data.TempalteCommon = tc
	data.Active = "/u/explorer"

	userID, _ := getSessionUserID(c)

	var list model.List
	model.DB.Where("id = ? and user = ?", id, userID).Preload("Epubs").Preload("Children").Find(&list)
	data.List = list

	if list.ID != 0 {
		return c.Render(http.StatusOK, "explorer", data)
	}

	return c.String(http.StatusNotFound, "not found")
}

type _deleteExplorerBody struct {
	Lists []string `json:"lists"`
	Books []string `json:"books"`
}

func deleteExplorer(c echo.Context) error {
	parentListID := c.Param("list_id")
	_body := new(_deleteExplorerBody)

	if err := c.Bind(_body); err != nil {
		return err
	}

	userID, _ := getSessionUserID(c)

	// remove child lists
	if len(_body.Lists) > 0 {
		model.DB.Where("parent_id = ? and user = ? and id in (?)", parentListID, userID, _body.Lists).Delete(model.List{})
	}

	// remove associated epubs
	if len(_body.Books) > 0 {
		var parentList model.List
		model.DB.Where("user = ? and id = ?", userID, parentListID).Find(&parentList)
		var epubs []model.Epub
		model.DB.Where("id in (?)", _body.Books).Find(&epubs)
		model.DB.Model(&parentList).Association("Epubs").Delete(epubs)
		model.DB.
			Where("user_id = ? and list_id = ? and epub_id in (?)",
				userID, parentListID, _body.Books).
			Delete(model.UserListEpub{})
	}
	return c.String(http.StatusOK, "")
}
