package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LogoutHandler clears the JWT token cookie and redirects to login
func LogoutHandler(c *gin.Context) {
	// Clear the token cookie
	c.SetCookie(
		"token",
		"",
		-1, // MaxAge -1 deletes the cookie
		"/",
		"",
		false,
		true, // HttpOnly
	)

	// Redirect to login page
	c.Redirect(http.StatusSeeOther, "/login")
}
