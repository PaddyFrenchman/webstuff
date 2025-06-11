package middleware

import (
	"go-web/internal/appcontext"

	"github.com/labstack/echo/v4"
)

func AppContext(cc *appcontext.AppContext) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc.Context = c
			return h(cc)
		}
	}
}
