package main

import (
	"database/sql"
	"log"

	"github.com/celest1al/simplebank-go/api"
	db "github.com/celest1al/simplebank-go/db/sqlc"
	"github.com/celest1al/simplebank-go/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	connection, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}

	store := db.NewStore(connection)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start the server", err)
	}
}
