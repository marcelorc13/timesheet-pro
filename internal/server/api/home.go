package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// Home example
// @Summary Home example
// @Schemes
// @Description returns welcome message
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Welcome!
// @Router / [get]
func HomeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "Welcome!")
}
