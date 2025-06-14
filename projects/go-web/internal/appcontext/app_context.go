package appcontext

import (
	"go-web/config"
	"go-web/internal/i18n"
	"go-web/internal/store"

	"github.com/labstack/echo/v4"
)

// AppContext is the new context in the request / response cycle
// We can use the db store, cache and central configuration
type AppContext struct {
	echo.Context
	UserStore store.User
	Cache     store.Cache
	Config    *config.Configuration
	Loc       i18n.I18ner
}
