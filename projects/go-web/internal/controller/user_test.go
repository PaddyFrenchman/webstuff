package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-web/internal/models"

	"github.com/labstack/echo/v4"
	"go-web/internal/appcontext"
	"go-web/internal/core/middleware"
	"github.com/stretchr/testify/assert"
)

type UserFakeStore struct{}

func (s *UserFakeStore) First(m *models.User) error {
	return nil
}
func (s *UserFakeStore) Find(m *[]models.User) error {
	return nil
}
func (s *UserFakeStore) Create(m *models.User) error {
	return nil
}
func (s *UserFakeStore) Ping() error {
	return nil
}

func TestUserPage(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/users/"+e.testUser.ID, nil)
	rec := httptest.NewRecorder()
	e.server.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUnitGetUserJson(t *testing.T) {
	s := echo.New()
	g := s.Group("/api")

	req := httptest.NewRequest(echo.GET, "/api/users/"+e.testUser.ID, nil)
	rec := httptest.NewRecorder()

	userCtrl := &User{}

	cc := &context.AppContext{
		Config:    e.config,
		UserStore: &UserFakeStore{},
	}

	s.Use(middleware.AppContext(cc))

	g.GET("/users/:id", userCtrl.GetUserJSON)
	s.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
