package mssql

import (
	"github.com/alpakih/go-api/pkg/database"
	"gorm.io/driver/sqlserver"
)

func init() {
	database.RegisterDialect("mssql", "sqlserver://{username}:{password}@{host}:{port}?database={name}&{options}", sqlserver.Open)
}
