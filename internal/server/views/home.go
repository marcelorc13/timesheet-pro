package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/marcelorc13/timesheet-pro/internal/domain"
	"github.com/marcelorc13/timesheet-pro/internal/repository"
	"github.com/marcelorc13/timesheet-pro/internal/utils"
)

type HomeViewHandler struct {
	orgRepo repository.OrganizationRepository
}

func NewHomeViewHandler(orgRepo repository.OrganizationRepository) *HomeViewHandler {
	return &HomeViewHandler{orgRepo: orgRepo}
}

func HomeHandler(c *gin.Context, orgRepo repository.OrganizationRepository) {
	// Check if user is logged in
	tokenString, err := c.Cookie("token")
	if err != nil {
		// No token, redirect to login
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Verify token and get user info
	claims, err := utils.GetTokenClaims(tokenString)
	if err != nil {
		// Invalid token, redirect to login
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Extract user ID
	userIDStr, ok := claims["id"].(string)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Get user's organization
	orgRes, err := orgRepo.GetOrganizationByUserID(c.Request.Context(), userID)
	
	// If user has no organization, redirect to create page
	if err != nil || !orgRes.Success {
		c.Redirect(http.StatusSeeOther, "/organizations/new")
		return
	}

	// User has an organization, get the organization data
	org, ok := orgRes.Data.(domain.Organization)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/organizations/new")
		return
	}

	// Redirect to organization detail page
	c.Redirect(http.StatusSeeOther, "/organizations/"+org.ID.String())
}
