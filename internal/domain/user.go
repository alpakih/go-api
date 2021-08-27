package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        string    `gorm:"column:id;type:varchar(60);primary_key:true" json:"id"`
	UserName  string    `gorm:"column:username;type:varchar(50);unique" json:"user_name"`
	Password  string    `gorm:"column:password;type:varchar(100)" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (c User) TableName() string {
	return "users"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *User) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()

	return
}

type TokenRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,max=100"`
}

type StoreRequest struct {
	Username string `json:"username" validate:"required,unique=username:users"`
	Password string `json:"password" validate:"required,max=100"`
}

type UpdateRequest struct {
	ID       string `json:"id" validate:"required"`
	Username string `json:"username" validate:"required,max=50,unique_update=ID:users:username:id"`
	Password string `json:"password" validate:"max=100"`
}

type UserService interface {
	Fetch(limit int, offset int) ([]User, error)
	GetByID(id string) (User, error)
	GetByUsername(username string) (User, error)
	Update(user UpdateRequest) error
	Store(user StoreRequest) error
	Delete(id string) error
}

type UserRepository interface {
	Fetch(limit int, offset int) ([]User, error)
	FindByID(id string) (User, error)
	FindByUsername(username string) (User, error)
	Update(user User) error
	Store(user User) error
	Delete(id string) error
}
