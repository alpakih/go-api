package service

import (
	"errors"
	"github.com/alpakih/go-api/internal/domain"
	"github.com/alpakih/go-api/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGetByID(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := domain.User{
		UserName: "testing",
		Password: "123123",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(mockUser, nil).Once()

		u := NewUserService(mockUserRepo)

		a, err := u.GetByID(mockUser.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockUserRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{}, errors.New("unexpected")).Once()

		u := NewUserService(mockUserRepo)

		a, err := u.GetByID(mockUser.ID)

		assert.Error(t, err)
		assert.Equal(t, domain.User{}, a)

		mockUserRepo.AssertExpectations(t)
	})

}
