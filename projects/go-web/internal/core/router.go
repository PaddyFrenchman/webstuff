package core

import (
	"go-web/internal/appcontext"
	mid "go-web/internal/core/middleware"
	"go-web/internal/i18n"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	v "gopkg.in/go-playground/validator.v9"
)

func NewRouter(server *Server) *echo.Echo {
	config := server.config
	e := echo.New()
	e.Validator = &Validator{validator: v.New()}

	cc := appcontext.AppContext{
		Cache:     &CacheStore{Cache: server.cache},
		Config:    config,
		UserStore: &UserStore{DB: server.db},
		Loc:       i18n.New(),
	}

	e.Use(mid.AppContext(&cc))

	if config.RequestLogger {
		e.Use(middleware.Logger()) // request logger
	}

	e.Use(middleware.Recover())       // panic errors are thrown
	e.Use(middleware.BodyLimit("5M")) // limit body payload to 5MB
	e.Use(middleware.Secure())        // provide protection against injection attacks
	e.Use(middleware.RequestID())     // generate unique requestId

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// add custom error formating
	e.HTTPErrorHandler = HTTPErrorHandler

	// Add html templates with go template syntax
	renderer := newTemplateRenderer(config.LayoutDir, config.TemplateDir)
	e.Renderer = renderer

	return e
}
