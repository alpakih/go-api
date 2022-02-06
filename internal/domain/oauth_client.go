package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type OauthClient struct {
	ID        string    `gorm:"column:id;type:varchar(60);primary_key:true" json:"id"`
	UserId    string    `gorm:"column:user_id;type:varchar(60)" json:"id"`
	ClientSecret  string `gorm:"column:secret;type:varchar(100)" json:"secret"`
	Provider  string    `gorm:"column:provider;type:varchar(100)" json:"provider"`
	Redirect  string    `gorm:"column:redirect;type:text" json:"redirect"`
	Revoked   int    `gorm:"column:revoked" json:"revoked"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"updated_at"`
}

func (c OauthClient) TableName() string {
	return "oauth_clients"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *OauthClient) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()

	return
}

type OauthTokenRequest struct {
	ID string `json:"client_id" validate:"required"`
	ClientSecret string `json:"client_secret" validate:"required"`
}
