package main

import (
	_ "github.com/lib/pq"
	"github.com/nanmenkaimak/report-management/internal/dbs/postgres"
	"github.com/nanmenkaimak/report-management/internal/handlers"
	"net/http"
)

const portNumber = ":8082"

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
