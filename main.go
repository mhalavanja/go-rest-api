package main

import (
	"database/sql"
	"log"

	"github.com/mhalavanja/go-rest-api/api"
	"github.com/mhalavanja/go-rest-api/db"
	"github.com/mhalavanja/go-rest-api/db/sqlc"
	"github.com/mhalavanja/go-rest-api/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configuration: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the database: ", err)
	}

	db.ExecuteStoredProcedures(conn)
	store := sqlc.New(conn)

	upgrader := api.NewUpgrader(config)
	hub := api.NewHub(&upgrader)
	go hub.Run()

	server, err := api.NewServer(config, store, hub)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
