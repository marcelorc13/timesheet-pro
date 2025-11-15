package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/marcelorc13/timesheet-pro/internal/repository"
	"github.com/marcelorc13/timesheet-pro/internal/server"
	"github.com/marcelorc13/timesheet-pro/internal/server/api"
	service "github.com/marcelorc13/timesheet-pro/internal/services"
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

	ur := repository.NewUserRepository(db)
	us := service.NewUserService(*ur)
	uh := api.NewUserHandler(*us)

	router.APIRoutes(*uh)
	router.ViewsRoutes()

	router.Start()
}
