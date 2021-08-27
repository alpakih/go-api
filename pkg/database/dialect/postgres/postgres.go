package postgres

import (
	"github.com/alpakih/go-api/pkg/database"
	"gorm.io/driver/postgres"
)

func init() {
	database.RegisterDialect("postgres", "host={host} port={port} user={username} dbname={name} password={password} {options}", postgres.Open)
}