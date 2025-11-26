package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marcelorc13/timesheet-pro/internal/templates/pages"
	"github.com/marcelorc13/timesheet-pro/internal/utils"
	service "github.com/marcelorc13/timesheet-pro/internal/services"
)

type ProfileViewHandler struct {
	userServ *service.UserService
}

func NewProfileViewHandler(userServ *service.UserService) *ProfileViewHandler {
	return &ProfileViewHandler{userServ: userServ}
}

// ProfilePageHandler shows the user profile page
func (h *ProfileViewHandler) ProfilePageHandler(c *gin.Context) {
	// Get user ID from JWT token
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	claims, err := utils.GetTokenClaims(tokenString)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	userIDStr, ok := claims["id"].(string)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	userName, ok := claims["name"].(string)
	if !ok {
		userName = ""
	}

	// Get user data
	user, err := h.userServ.GetByID(c.Request.Context(), userIDStr)
	if err != nil || user == nil {
		c.String(http.StatusNotFound, "Usuário não encontrado")
		return
	}

	utils.Render(c.Request.Context(), c.Writer, pages.ProfilePage(*user, userName))
}
