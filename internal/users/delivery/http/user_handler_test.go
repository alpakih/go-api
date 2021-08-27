package http

import (
	"github.com/alpakih/go-api/internal/domain"
	"github.com/alpakih/go-api/internal/domain/mocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetByID(t *testing.T) {
	id := uuid.New().String()
	mockUser := domain.User{
		ID:        id,
		UserName:  "testing",
		Password:  "password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUCase := new(mocks.UserService)

	mockUCase.On("GetByID", id).Return(mockUser, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/api/v1/users/"+id, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("api/v1/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(id)
	handler := UserHandler{
		UserService: mockUCase,
	}
	err = handler.GetUserByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}
