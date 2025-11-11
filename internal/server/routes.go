package server

import (
	"github.com/marcelorc13/timesheet-pro/docs"
	"github.com/marcelorc13/timesheet-pro/internal/server/api"
	"github.com/marcelorc13/timesheet-pro/internal/server/views"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (r Router) APIRoutes() {
   	docs.SwaggerInfo.BasePath = "/api/v1"
	apiRouter := r.Router.Group("/api/v1")
	apiRouter.GET("/", api.HomeHandler)

	r.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// http://localhost:port/swagger/index.html
}

func (r Router) ViewsRoutes() {
	viewsRouter := r.Router.Group("/")
	viewsRouter.GET("/", views.HomeHandler)
}
