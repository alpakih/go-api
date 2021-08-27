package mysql

import (
	"github.com/alpakih/go-api/pkg/database"
	"gorm.io/driver/mysql"
)

func init() {
	database.RegisterDialect("mysql", "{username}:{password}@({host}:{port})/{name}?{options}", mysql.Open)
}
