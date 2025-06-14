package controller

import (
	"net/http"

	"go-web/internal/appcontext"

	"github.com/labstack/echo/v4"
)

type Healthcheck struct{}

type healthcheckReport struct {
	Health  string          `json:"health"`
	Details map[string]bool `json:"details"`
}

// GetHealthcheck returns the current functional state of the application
func (ctrl Healthcheck) GetHealthcheck(c echo.Context) error {
	cc := c.(*appcontext.AppContext)
	m := healthcheckReport{Health: "OK"}

	dbCheck := cc.UserStore.Ping()
	cacheCheck := cc.Cache.Ping()

	if dbCheck != nil {
		m.Health = "NOT"
		m.Details["db"] = false
	}

	if cacheCheck != nil {
		m.Health = "NOT"
		m.Details["cache"] = false
	}

	return c.JSON(http.StatusOK, m)
}
