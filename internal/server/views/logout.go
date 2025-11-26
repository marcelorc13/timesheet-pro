package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogoutHandler(c *gin.Context) {
	// Clear the token cookie with same settings as login
	c.SetCookie(
		"token",
		"",
		-1, // MaxAge -1 deletes the cookie
		"/",
		"", // Empty domain works on any domain
		true, // Secure flag for HTTPS
		true, // HttpOnly
	)

	// Redirect to login page
	c.Redirect(http.StatusSeeOther, "/login")
}
