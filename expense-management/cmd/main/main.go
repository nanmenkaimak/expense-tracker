package main

import (
	_ "github.com/lib/pq"
	"github.com/nanmenkaimak/expense-management/internal/dbs/postgres"
	"github.com/nanmenkaimak/expense-management/internal/handlers"
	"net/http"
)

const portNumber = ":8081"

func main() {
	db, err := postgres.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := handlers.NewRepo(db)
	handlers.NewHandlers(repo)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}

	err = srv.ListenAndServe()
}
