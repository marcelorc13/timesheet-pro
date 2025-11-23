package server

import (
	"github.com/marcelorc13/timesheet-pro/docs"
	"github.com/marcelorc13/timesheet-pro/internal/server/api"
	"github.com/marcelorc13/timesheet-pro/internal/server/views"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (r Router) APIRoutes(uh api.UserHandler, oh api.OrganizationHandler) {
	docs.SwaggerInfo.BasePath = "/api/v1"

	apiRouter := r.Router.Group("/api/v1")
	apiRouter.GET("/", api.HomeHandler)

	userRoutes := apiRouter.Group("users/")

	userRoutes.GET("/", uh.List)
	userRoutes.GET("/:id", uh.GetByID)
	userRoutes.DELETE("/:id", uh.Delete)
	userRoutes.POST("/", uh.Create)
	userRoutes.POST("/login", uh.Login)


	organizationRoutes := apiRouter.Group("organizations/")

	organizationRoutes.POST("/", oh.Create)

	r.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// http://localhost:port/swagger/index.html
}

func (r Router) ViewsRoutes() {
	viewsRouter := r.Router.Group("/")

	viewsRouter.GET("/signup", views.SignupHandler)
	viewsRouter.GET("/login", views.LoginHandler)

	authRoutes := viewsRouter.Group("/")
	authRoutes.Use(AuthMiddleware())
	authRoutes.GET("/", views.HomeHandler)
}
