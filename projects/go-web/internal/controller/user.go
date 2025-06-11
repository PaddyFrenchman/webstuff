package controller

import (
	"go-web/internal/appcontext"
	"net/http"

	"go-web/internal/core/errors"
	"go-web/internal/models"

	"github.com/labstack/echo/v4"
)

type (
	User          struct{}
	UserViewModel struct {
		Name string
		ID   string
	}
)

func (ctrl User) GetUser(c echo.Context) error {
	cc := c.(*appcontext.AppContext)
	userID := c.Param("id")

	user := models.User{ID: userID}

	err := cc.UserStore.First(&user)

	if err != nil {
		b := errors.NewBoom(errors.UserNotFound, errors.ErrorText(errors.UserNotFound), err)
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, b)
	}

	vm := UserViewModel{
		Name: user.Name,
		ID:   user.ID,
	}

	return c.Render(http.StatusOK, "user.html", vm)

}

func (ctrl User) GetUserJSON(c echo.Context) error {
	cc := c.(*appcontext.AppContext)
	userID := c.Param("id")

	user := models.User{ID: userID}

	err := cc.UserStore.First(&user)

	if err != nil {
		b := errors.NewBoom(errors.UserNotFound, errors.ErrorText(errors.UserNotFound), err)
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, b)
	}

	vm := UserViewModel{
		Name: user.Name,
		ID:   user.ID,
	}

	return c.JSON(http.StatusOK, vm)
}
