package dbrepo

import (
	"database/sql"

	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/config"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}
type testDbRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDbRepo{
		App: a,
	}
}

func NewpostgresDBRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
