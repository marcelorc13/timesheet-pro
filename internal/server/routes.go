package server

import (
	"github.com/marcelorc13/timesheet-pro/internal/server/api"
	"github.com/marcelorc13/timesheet-pro/internal/server/views"
)


func (r Router) APIRoutes() {
	apiRouter := r.Router.Group("/api/v1")
	apiRouter.GET("/", api.HomeHandler)

	viewsRouter := r.Router.Group("/")
	viewsRouter.GET("/", views.HomeHandler)
}
