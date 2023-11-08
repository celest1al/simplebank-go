package main

import (
	"database/sql"
	"log"

	"github.com/celest1al/simplebank-go/api"
	db "github.com/celest1al/simplebank-go/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver     = "postgres"
	dbSource     = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddres = "0.0.0.0:8080"
)

func main() {
	connection, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}

	store := db.NewStore(connection)
	server := api.NewServer(store)

	err = server.Start(serverAddres)

	if err != nil {
		log.Fatal("Cannot start the server", err)
	}
}
