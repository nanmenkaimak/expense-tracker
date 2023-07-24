package dbrepo

import (
	"github.com/jmoiron/sqlx"
	"github.com/nanmenkaimak/expense-management/internal/repository"
)

type postgresDBRepo struct {
	DB *sqlx.DB
}

func NewPostgresRepo(conn *sqlx.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		DB: conn,
	}
}
