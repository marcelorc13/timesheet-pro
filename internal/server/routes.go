package server

import (
	"github.com/gin-gonic/gin"
	"github.com/marcelorc13/timesheet-pro/docs"
	"github.com/marcelorc13/timesheet-pro/internal/repository"
	"github.com/marcelorc13/timesheet-pro/internal/server/api"
	"github.com/marcelorc13/timesheet-pro/internal/server/views"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (r Router) APIRoutes(uh api.UserHandler, oh api.OrganizationHandler, th api.TimesheetHandler) {
	docs.SwaggerInfo.BasePath = "/api/v1"

	apiRouter := r.Router.Group("/api/v1")
	apiRouter.GET("/", api.HomeHandler)

	userRoutes := apiRouter.Group("users/")

	userRoutes.GET("/", uh.List)
	userRoutes.GET("/:id", uh.GetByID)
	userRoutes.DELETE("/:id", uh.Delete)
	userRoutes.PUT("/:id", uh.UpdateMyProfile)
	userRoutes.POST("/", uh.Create)
	userRoutes.POST("/login", uh.Login)


	organizationRoutes := apiRouter.Group("organizations/")

	organizationRoutes.POST("/", oh.Create)
	organizationRoutes.GET("/", oh.List)
	organizationRoutes.GET("/:id", oh.GetByID)
	organizationRoutes.GET("/user/:userId", oh.GetByUserID)
	organizationRoutes.PUT("/:id", oh.Update)
	organizationRoutes.DELETE("/:id", oh.Delete)
	organizationRoutes.POST("/:id/users", oh.AddUser)
	organizationRoutes.DELETE("/:id/users/:userId", oh.RemoveUser)
	organizationRoutes.POST("/:id/leave", oh.Leave)

	organizationRoutes.POST("/:id/clock-in", th.ClockIn)
	organizationRoutes.GET("/:id/timesheets/me", th.GetMyTimesheets)
	organizationRoutes.GET("/:id/timesheets/me/status", th.GetMyStatus)
	organizationRoutes.GET("/:id/users/:userId/timesheets", th.GetUserTimesheets)
	organizationRoutes.GET("/:id/timesheets/all", th.GetAllTimesheets)

	r.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// http://localhost:port/swagger/index.html
}

func (r Router) ViewsRoutes(ovh views.OrganizationViewHandler, tvh views.TimesheetViewHandler, pvh views.ProfileViewHandler, orgRepo *repository.OrganizationRepository) {
	viewsRouter := r.Router.Group("/")

	viewsRouter.GET("/signup", views.SignupHandler)
	viewsRouter.GET("/login", views.LoginHandler)
	viewsRouter.GET("/logout", views.LogoutHandler)

	authRoutes := viewsRouter.Group("/")
	authRoutes.Use(AuthMiddleware())
	authRoutes.GET("/", func(c *gin.Context) {
		views.HomeHandler(c, *orgRepo)
	})
	
	authRoutes.GET("/organizations/new", ovh.OrganizationCreateHandler)
	authRoutes.GET("/organizations/:id", ovh.OrganizationDetailHandler)
	authRoutes.GET("/organizations/:id/edit", ovh.OrganizationEditHandler)
	authRoutes.GET("/organizations/:id/add-user", ovh.OrganizationAddUserHandler)
	
	authRoutes.GET("/timesheet", tvh.TimesheetPageHandler)

	authRoutes.GET("/profile", pvh.ProfilePageHandler)
}
