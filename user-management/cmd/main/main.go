package main

import (
	_ "github.com/lib/pq"
	"github.com/nanmenkaimak/user-management/internal/dbs/postgres"
	"github.com/nanmenkaimak/user-management/internal/handlers"
	"net/http"
)

const portNumber = ":8080"

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
