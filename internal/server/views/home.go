package views

import (
	"github.com/gin-gonic/gin"
	"github.com/marcelorc13/timesheet-pro/internal/templates/pages"
	"github.com/marcelorc13/timesheet-pro/internal/utils"
)

func HomeHandler(c *gin.Context) {
	utils.Render(c.Request.Context(), c.Writer, pages.HomePage())
}
