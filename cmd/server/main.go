package main

import (
	"log"

	"github.com/harry713j/vibe_writer/internal/config"
	"github.com/harry713j/vibe_writer/internal/db"
	"github.com/harry713j/vibe_writer/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	dbConfig := config.LoadDBConfig()
	db := db.ConnectDB(dbConfig)

	defer db.Close()

	serverConfig := server.LoadServerConfig()
	srv := server.NewServer(serverConfig)

	log.Println("Server has started on Port: ", serverConfig.Port)
	log.Fatal(srv.ListenAndServe())
}
