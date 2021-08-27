package mysql

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	id := uuid.New().String()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})


	mock.ExpectQuery(
		`SELECT(.*)`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at"}).
			AddRow(id, "testing", "password", time.Now(), time.Now()))


	a := NewMysqlUserRepository(gormDB)

	anUser, err := a.FindByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, anUser)
}
