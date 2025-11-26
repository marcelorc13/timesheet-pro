// Package api provides HTTP handlers for REST API endpoints.
// It handles request parsing, validation, and response formatting for the application.
package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/marcelorc13/timesheet-pro/internal/domain"
	service "github.com/marcelorc13/timesheet-pro/internal/services"
	"github.com/marcelorc13/timesheet-pro/internal/utils"
)

type TimesheetHandler struct {
	service *service.TimesheetService
}

func NewTimesheetHandler(ts *service.TimesheetService) *TimesheetHandler {
	return &TimesheetHandler{ts}
}

// ClockIn handles POST /api/v1/organizations/:id/clock-in
// Clocks user in or out based on current status
func (h *TimesheetHandler) ClockIn(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID da organização inválido"})
		return
	}

	// Get user ID from JWT token
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token não encontrado"})
		return
	}

	claims, err := utils.GetTokenClaims(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token inválido"})
		return
	}

	userIDStr, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token inválido"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "ID de usuário inválido no token"})
		return
	}

	// Call service
	err = h.service.ClockIn(c.Request.Context(), userID, orgID)
	if err != nil {
		// Log the error for debugging
		println("ClockIn error:", err.Error())
		
		if err.Error() == "usuário não é membro desta organização" {
			c.JSON(http.StatusForbidden, domain.HttpResponse{Status: http.StatusForbidden, Message: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Registro de ponto realizado com sucesso"})
}

// GetMyTimesheets handles GET /api/v1/organizations/:id/timesheets/me
// Returns authenticated user's timesheets for date range (defaults to today)
func (h *TimesheetHandler) GetMyTimesheets(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID da organização inválido"})
		return
	}

	// Get user ID from JWT token
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token não encontrado"})
		return
	}

	claims, err := utils.GetTokenClaims(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token inválido"})
		return
	}

	userIDStr, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token inválido"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "ID de usuário inválido no token"})
		return
	}

	// Parse query parameters for date range
	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr != "" && endStr != "" {
		// Parse dates
		startDate, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "Formato de data inicial inválido (use YYYY-MM-DD)"})
			return
		}
		endDate, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "Formato de data final inválido (use YYYY-MM-DD)"})
			return
		}

		// Get timesheets for date range
		timesheets, err := h.service.GetUserTimesheets(c.Request.Context(), userID, userID, orgID, startDate, endDate)
		if err != nil {
			if err.Error() == "usuário não é membro desta organização" {
				c.JSON(http.StatusForbidden, domain.HttpResponse{Status: http.StatusForbidden, Message: err.Error()})
				return
			}
			c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Timesheets encontrados", Data: timesheets})
		return
	}

	// Default to today
	today := time.Now().Truncate(24 * time.Hour)
	timesheet, err := h.service.GetUserTimesheet(c.Request.Context(), userID, orgID, today)
	if err != nil {
		if err.Error() == "timesheet não encontrado para esta data" {
			c.JSON(http.StatusNotFound, domain.HttpResponse{Status: http.StatusNotFound, Message: err.Error()})
			return
		}
		if err.Error() == "usuário não é membro desta organização" {
			c.JSON(http.StatusForbidden, domain.HttpResponse{Status: http.StatusForbidden, Message: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Timesheet encontrado", Data: timesheet})
}

// GetMyStatus handles GET /api/v1/organizations/:id/timesheets/me/status
// Returns current clock in/out status
func (h *TimesheetHandler) GetMyStatus(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID da organização inválido"})
		return
	}

	// Get user ID from JWT token
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token não encontrado"})
		return
	}

	claims, err := utils.GetTokenClaims(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token inválido"})
		return
	}

	userIDStr, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token inválido"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "ID de usuário inválido no token"})
		return
	}

	// Get status
	status, timestamp, err := h.service.GetCurrentStatus(c.Request.Context(), userID, orgID)
	if err != nil {
		if err.Error() == "usuário não é membro desta organização" {
			c.JSON(http.StatusForbidden, domain.HttpResponse{Status: http.StatusForbidden, Message: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	response := map[string]interface{}{
		"status":    status,
		"timestamp": timestamp,
	}

	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Status atual", Data: response})
}

// GetUserTimesheets handles GET /api/v1/organizations/:id/users/:userId/timesheets
// Admin only - returns specific user's timesheets
func (h *TimesheetHandler) GetUserTimesheets(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID da organização inválido"})
		return
	}

	targetUserIDStr := c.Param("userId")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID do usuário inválido"})
		return
	}

	// Get requesting user ID from JWT token
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token não encontrado"})
		return
	}

	claims, err := utils.GetTokenClaims(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token inválido"})
		return
	}

	requestingUserIDStr, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token inválido"})
		return
	}

	requestingUserID, err := uuid.Parse(requestingUserIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "ID de usuário inválido no token"})
		return
	}

	// Parse date range (defaults to last 30 days)
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	if startParam := c.Query("start"); startParam != "" {
		parsedStart, err := time.Parse("2006-01-02", startParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "Formato de data inicial inválido (use YYYY-MM-DD)"})
			return
		}
		startDate = parsedStart
	}

	if endParam := c.Query("end"); endParam != "" {
		parsedEnd, err := time.Parse("2006-01-02", endParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "Formato de data final inválido (use YYYY-MM-DD)"})
			return
		}
		endDate = parsedEnd
	}

	// Call service
	timesheets, err := h.service.GetUserTimesheets(c.Request.Context(), requestingUserID, targetUserID, orgID, startDate, endDate)
	if err != nil {
		if err.Error() == "apenas administradores podem visualizar timesheets de outros usuários" {
			c.JSON(http.StatusForbidden, domain.HttpResponse{Status: http.StatusForbidden, Message: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Data: timesheets})
}

// GetAllTimesheets handles GET /api/v1/organizations/:id/timesheets/all
// Admin only - returns all organization timesheets for a date
func (h *TimesheetHandler) GetAllTimesheets(c *gin.Context) {
	orgIDStr := c.Param("id")
	_ = orgIDStr
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID da organização inválido"})
		return
	}

	// Get requesting user ID from JWT token
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token não encontrado"})
		return
	}

	claims, err := utils.GetTokenClaims(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token inválido"})
		return
	}

	adminUserIDStr, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token inválido"})
		return
	}

	adminUserID, err := uuid.Parse(adminUserIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "ID de usuário inválido no token"})
		return
	}

	// Parse date (defaults to today)
	date := time.Now().Truncate(24 * time.Hour)
	if dateParam := c.Query("date"); dateParam != "" {
		parsedDate, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "Formato de data inválido (use YYYY-MM-DD)"})
			return
		}
		date = parsedDate
	}

	// Get all timesheets
	timesheets, err := h.service.GetOrganizationTimesheets(c.Request.Context(), adminUserID, orgID, date)
	if err != nil {
		if err.Error() == "apenas administradores podem visualizar timesheets da organização" {
			c.JSON(http.StatusForbidden, domain.HttpResponse{Status: http.StatusForbidden, Message: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Timesheets da organização", Data: timesheets})
}

// GetTimesheetByID handles GET /api/v1/timesheets/:id
// Returns a single timesheet by ID (user must own it or be admin)
func (h *TimesheetHandler) GetTimesheetByID(c *gin.Context) {
	timesheetIDStr := c.Param("id")
	timesheetID, err := uuid.Parse(timesheetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID do timesheet inválido"})
		return
	}

	// Get user ID from JWT token
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token não encontrado"})
		return
	}

	claims, err := utils.GetTokenClaims(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token inválido"})
		return
	}

	userIDStr, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token inválido"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "ID de usuário inválido no token"})
		return
	}

	// Call service
	timesheet, err := h.service.GetTimesheetByID(c.Request.Context(), userID, timesheetID)
	if err != nil {
		if err.Error() == "você não tem permissão para visualizar este timesheet" {
			c.JSON(http.StatusForbidden, domain.HttpResponse{Status: http.StatusForbidden, Message: err.Error()})
			return
		}
		if err.Error() == "timesheet não encontrado" {
			c.JSON(http.StatusNotFound, domain.HttpResponse{Status: http.StatusNotFound, Message: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Data: timesheet})
}
