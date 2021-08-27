package service

import (
	"github.com/alpakih/go-api/internal/domain"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository domain.UserRepository
}

// NewUserService will create new an userService object representation of domain.UserService interface
func NewUserService(ur domain.UserRepository) domain.UserService {
	return &userService{
		userRepository: ur,
	}
}

func (u userService) Fetch(limit int, offset int) ([]domain.User, error) {
	return u.userRepository.Fetch(limit, offset)
}

func (u userService) GetByID(id string) (domain.User, error) {
	return u.userRepository.FindByID(id)
}

func (u userService) Update(user domain.UpdateRequest) error {
	var entity domain.User

	entity.ID = user.ID
	entity.UserName = user.Username
	if user.Password != "" {
		entity.Password = user.Password
	}
	return u.userRepository.Update(entity)

}

func (u userService) Store(user domain.StoreRequest) error {
	if bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), viper.GetInt("app.bcryptCost"));
	err != nil {
		return err
	} else {
		entity := domain.User{
			UserName: user.Username,
			Password: string(bytes),
		}
		return u.userRepository.Store(entity)
	}
}

func (u userService) Delete(id string) error {
	return u.userRepository.Delete(id)
}

func (u userService) GetByUsername(id string) (domain.User, error) {
	return u.userRepository.FindByUsername(id)
}