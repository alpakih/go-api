package mysql

import (
	"github.com/alpakih/go-api/internal/domain"
	"github.com/alpakih/go-api/pkg/database"
	"gorm.io/gorm"
)

type mysqlUserRepo struct {
	DB *gorm.DB
}

// NewMysqlUserRepository will create an implementation of domain.UserRepository
func NewMysqlUserRepository(db *gorm.DB) domain.UserRepository {
	return &mysqlUserRepo{
		DB: db,
	}
}

func (m mysqlUserRepo) Fetch(limit int, offset int) ([]domain.User, error) {
	var entity []domain.User
	paginator := database.NewPaginator(m.DB, offset, limit, &entity)
	if err := paginator.Find().Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (m mysqlUserRepo) FindByID(id string) (domain.User, error) {
	var entity domain.User
	if err := m.DB.First(&entity, "id =?", id).Error; err != nil {
		return domain.User{}, err
	}
	return entity, nil
}

func (m mysqlUserRepo) Update(user domain.User) error {
	return m.DB.Updates(&user).Error
}

func (m mysqlUserRepo) Store(user domain.User) error {
	return m.DB.Create(&user).Error
}

func (m mysqlUserRepo) Delete(id string) error {
	return m.DB.Delete(&domain.User{}, "id =?", id).Error
}

func (m mysqlUserRepo) FindByUsername(username string) (domain.User, error) {
	var entity domain.User
	if err := m.DB.First(&entity, "username =?", username).Error; err != nil {
		return domain.User{}, err
	}
	return entity, nil
}

