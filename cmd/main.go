package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/marcelorc13/timesheet-pro/internal/repository"
	"github.com/marcelorc13/timesheet-pro/internal/server"
)

func main() {
	_ = godotenv.Load()
	connString := os.Getenv("POSTGRES_URL")

	ctx := context.Background()

	db := repository.NewPool(ctx, connString)
	if err := db.Ping(ctx); err != nil {
		panic(err)
	}

	r := gin.Default()

	router := server.NewRouter(r)

	router.APIRoutes()

	router.Start()
}
