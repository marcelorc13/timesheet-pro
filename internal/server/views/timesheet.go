// Package views provides HTTP handlers for rendering HTML pages using templ templates.
// It handles authentication, authorization, and renders server-side HTML for the web interface.
package views

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	service "github.com/marcelorc13/timesheet-pro/internal/services"
	"github.com/marcelorc13/timesheet-pro/internal/templates/pages"
	"github.com/marcelorc13/timesheet-pro/internal/utils"
	"github.com/marcelorc13/timesheet-pro/internal/domain"
)

type TimesheetViewHandler struct {
	timesheetServ *service.TimesheetService
	orgServ       *service.OrganizationService
}

func NewTimesheetViewHandler(timesheetServ *service.TimesheetService, orgServ *service.OrganizationService) *TimesheetViewHandler {
	return &TimesheetViewHandler{
		timesheetServ: timesheetServ,
		orgServ:       orgServ,
	}
}

// TimesheetPageHandler shows the user's timesheet page with clock in/out
func (h *TimesheetViewHandler) TimesheetPageHandler(c *gin.Context) {
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

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Get user name
	userName, ok := claims["name"].(string)
	if !ok {
		userName = ""
	}

	// Get user's organization
	org, err := h.orgServ.GetOrganizationByUserID(c.Request.Context(), userID)
	if err != nil || org == nil {
		c.String(http.StatusNotFound, "Você não pertence a nenhuma organização")
		return
	}

	// Get today's timesheet
	today := time.Now().Truncate(24 * time.Hour)
	timesheet, _ := h.timesheetServ.GetUserTimesheet(c.Request.Context(), userID, org.ID, today)

	// Get current status
	status, timestamp, _ := h.timesheetServ.GetCurrentStatus(c.Request.Context(), userID, org.ID)
	
	var lastTimestampStr *string
	if timestamp != nil {
		formatted := timestamp.Format("15:04:05")
		lastTimestampStr = &formatted
	}

	utils.Render(c.Request.Context(), c.Writer, pages.TimesheetPage(*org, timesheet, status, lastTimestampStr, userName))
}

// AdminTimesheetPageHandler shows admin view of all organization timesheets
func (h *TimesheetViewHandler) AdminTimesheetPageHandler(c *gin.Context) {
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

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Get user name
	userName, ok := claims["name"].(string)
	if !ok {
		userName = ""
	}

	// Get user's organization
	org, err := h.orgServ.GetOrganizationByUserID(c.Request.Context(), userID)
	if err != nil || org == nil {
		c.String(http.StatusNotFound, "Você não pertence a nenhuma organização")
		return
	}

	// Check if user is admin
	isAdmin, err := h.orgServ.IsUserAdmin(c.Request.Context(), userID, org.ID)
	if err != nil || !isAdmin {
		c.String(http.StatusForbidden, "Apenas administradores podem acessar esta página")
		return
	}

	// Get all timesheets for today
	today := time.Now().Truncate(24 * time.Hour)
	timesheets, err := h.timesheetServ.GetOrganizationTimesheets(c.Request.Context(), userID, org.ID, today)
	if err != nil {
		timesheets = []domain.UserTimesheet{}
	}

	utils.Render(c.Request.Context(), c.Writer, pages.AdminTimesheetPage(*org, timesheets, userName))
}
