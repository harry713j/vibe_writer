package main

import (
	"log"
	"os"
	"time"

	"github.com/harry713j/vibe_writer/internal/app"
	"github.com/harry713j/vibe_writer/internal/config"
	"github.com/harry713j/vibe_writer/internal/db"
	"github.com/harry713j/vibe_writer/internal/handler"
	"github.com/harry713j/vibe_writer/internal/repo"
	"github.com/harry713j/vibe_writer/internal/server"
	"github.com/harry713j/vibe_writer/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	// load the environment variables
	godotenv.Load()
	dbConfig := config.LoadDBConfig()
	db, err := db.ConnectDB(dbConfig)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	serverConfig := server.LoadServerConfig()

	userRepo := repo.NewUserRepository(db)
	refreshTokenRepo := repo.NewRefreshTokenRepository(db)
	jwtSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	accessTokenTTL := 15 * time.Minute

	authService := service.NewAuthService(userRepo, refreshTokenRepo, jwtSecret, accessTokenTTL)

	app := &app.App{
		AuthService: authService,
		AuthHandler: handler.NewAuthHandler(authService),
	}

	srv := server.NewServer(serverConfig, app)

	log.Println("Server has started on Port: ", serverConfig.Port)
	log.Fatal(srv.ListenAndServe())
}
