package route

import (
	"fmt"
	"net/http"

	"github.com/kyicy/readimension/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func getExplorerRoot(c echo.Context) error {
	return c.Redirect(http.StatusFound, fmt.Sprintf("/u/explorer/%v", 1))
}

type getBooksData struct {
	*TemplateCommon
	List    model.List
	HasUser bool
}

func getExplorer(c echo.Context) error {
	id := c.Param("list_id")

	tc := newTemplateCommon(c, "Library Explorer")
	data := &getBooksData{}
	data.TemplateCommon = tc
	data.Active = "/u/explorer"

	userID, _ := getSessionUserID(c)
	data.HasUser = (userID != "")

	var list model.List
	model.DB.Where("id = ?", id).Preload("Epubs", func(db *gorm.DB) *gorm.DB {
		return db.Order("epubs.title asc")
	}).Preload("Children").Find(&list)
	data.List = list

	if list.ID != 0 {
		return c.Render(http.StatusOK, "explorer", data)
	}

	return c.Redirect(http.StatusFound, "/sign-up")
}

type deleteExplorerBody struct {
	Lists []string `json:"lists"`
	Books []string `json:"books"`
}

func deleteExplorer(c echo.Context) error {
	parentListID := c.Param("list_id")
	body := new(deleteExplorerBody)

	if err := c.Bind(body); err != nil {
		return err
	}

	userID, _ := getSessionUserID(c)

	// remove child lists
	if len(body.Lists) > 0 {
		model.DB.Where("parent_id = ? and user = ? and id in (?)", parentListID, userID, body.Lists).Delete(model.List{})
	}

	// remove associated epubs
	if len(body.Books) > 0 {
		var parentList model.List
		model.DB.Where("user = ? and id = ?", userID, parentListID).Find(&parentList)
		var epubs []model.Epub
		model.DB.Where("id in (?)", body.Books).Find(&epubs)
		model.DB.Model(&parentList).Association("Epubs").Delete(&epubs)
		model.DB.
			Where("user_id = ? and list_id = ? and epub_id in (?)",
				userID, parentListID, body.Books).
			Delete(model.UserListEpub{})
	}
	return c.String(http.StatusOK, "")
}
