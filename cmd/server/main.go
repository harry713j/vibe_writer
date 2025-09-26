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
	err := godotenv.Load()

	if err != nil {
		log.Println(".env file not found")
	}

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
	profileRepo := repo.NewUserProfileRepository(db)
	blogRepo := repo.NewBlogRepository(db)
	commentRepo := repo.NewCommentRepository(db)

	authService := service.NewAuthService(userRepo, profileRepo, refreshTokenRepo, jwtSecret, accessTokenTTL)
	userProfileService := service.NewUserProfileService(profileRepo, userRepo)
	blogService := service.NewBlogService(blogRepo, userRepo)
	commentService := service.NewCommentService(commentRepo, blogRepo, userRepo)

	app := &app.App{
		AuthService:        authService,
		UserProfileService: userProfileService,
		BlogService:        blogService,
		CommentService:     commentService,

		AuthHandler:        handler.NewAuthHandler(authService),
		UserProfileHandler: handler.NewUserProfileHandler(userProfileService),
		BlogHandler:        handler.NewBlogHandler(blogService),
		CommentHandler:     handler.NewCommentHandler(commentService),
	}

	srv := server.NewServer(serverConfig, app)

	log.Println("Server has started on Port: ", serverConfig.Port)
	log.Fatal(srv.ListenAndServe())
}
