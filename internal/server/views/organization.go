// Package views provides HTTP handlers for rendering HTML pages using templ templates.
// It handles authentication, authorization, and renders server-side HTML for the web interface.
package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/marcelorc13/timesheet-pro/internal/domain"
	"github.com/marcelorc13/timesheet-pro/internal/repository"
	"github.com/marcelorc13/timesheet-pro/internal/templates/pages"
	"github.com/marcelorc13/timesheet-pro/internal/utils"
)

type OrganizationViewHandler struct {
	orgRepo  repository.OrganizationRepository
	userRepo repository.UserRepository
}

func NewOrganizationViewHandler(orgRepo repository.OrganizationRepository, userRepo repository.UserRepository) *OrganizationViewHandler {
	return &OrganizationViewHandler{
		orgRepo:  orgRepo,
		userRepo: userRepo,
	}
}

// OrganizationDetailHandler shows the details of a specific organization
func (h *OrganizationViewHandler) OrganizationDetailHandler(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "ID inválido")
		return
	}

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

	// Get organization
	res, err := h.orgRepo.GetByID(c.Request.Context(), orgID)
	if err != nil || !res.Success {
		c.String(http.StatusNotFound, "Organização não encontrada")
		return
	}

	org, ok := res.Data.(domain.Organization)
	if !ok {
		c.String(http.StatusInternalServerError, "Erro ao processar dados")
		return
	}

	// Check if user is admin
	adminRes, _ := h.orgRepo.IsUserAdmin(c.Request.Context(), userID, orgID)
	isAdmin := false
	if adminRes.Success {
		isAdmin, _ = adminRes.Data.(bool)
	}

	// Get organization members - initialize empty slice
	members := []domain.OrganizationUser{}
	membersRes, err := h.orgRepo.GetOrganizationMembers(c.Request.Context(), orgID)
	if err == nil && membersRes.Success && membersRes.Data != nil {
		// Try to cast to the correct type
		if membersList, ok := membersRes.Data.([]domain.OrganizationUser); ok {
			members = membersList
		}
	}

	// Get user name
	userName, ok := claims["name"].(string)
	if !ok {
		userName = ""
	}

	utils.Render(c.Request.Context(), c.Writer, pages.OrganizationDetailPage(org, isAdmin, userIDStr, members, userName))
}

// OrganizationCreateHandler shows the create organization form
func (h *OrganizationViewHandler) OrganizationCreateHandler(c *gin.Context) {
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

	// Get user name
	userName, ok := claims["name"].(string)
	if !ok {
		userName = ""
	}

	utils.Render(c.Request.Context(), c.Writer, pages.OrganizationCreatePage(userIDStr, userName))
}

// OrganizationEditHandler shows the edit organization form
func (h *OrganizationViewHandler) OrganizationEditHandler(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "ID inválido")
		return
	}

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

	// Check if user is admin
	adminRes, _ := h.orgRepo.IsUserAdmin(c.Request.Context(), userID, orgID)
	isAdmin := false
	if adminRes.Success {
		isAdmin, _ = adminRes.Data.(bool)
	}

	if !isAdmin {
		c.String(http.StatusForbidden, "Apenas administradores podem editar a organização")
		return
	}

	// Get organization
	res, err := h.orgRepo.GetByID(c.Request.Context(), orgID)
	if err != nil || !res.Success {
		c.String(http.StatusNotFound, "Organização não encontrada")
		return
	}

	org, ok := res.Data.(domain.Organization)
	if !ok {
		c.String(http.StatusInternalServerError, "Erro ao processar dados")
		return
	}

	// Get user name
	userName, ok := claims["name"].(string)
	if !ok {
		userName = ""
	}

	utils.Render(c.Request.Context(), c.Writer, pages.OrganizationEditPage(org, userName))
}

// OrganizationAddUserHandler shows the add user form
func (h *OrganizationViewHandler) OrganizationAddUserHandler(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "ID inválido")
		return
	}

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

	// Check if user is admin
	adminRes, _ := h.orgRepo.IsUserAdmin(c.Request.Context(), userID, orgID)
	isAdmin := false
	if adminRes.Success {
		isAdmin, _ = adminRes.Data.(bool)
	}

	if !isAdmin {
		c.String(http.StatusForbidden, "Apenas administradores podem adicionar usuários")
		return
	}

	// Get organization
	res, err := h.orgRepo.GetByID(c.Request.Context(), orgID)
	if err != nil || !res.Success {
		c.String(http.StatusNotFound, "Organização não encontrada")
		return
	}

	org, ok := res.Data.(domain.Organization)
	if !ok {
		c.String(http.StatusInternalServerError, "Erro ao processar dados")
		return
	}

	// Get user name
	userName, ok := claims["name"].(string)
	if !ok {
		userName = ""
	}

	utils.Render(c.Request.Context(), c.Writer, pages.OrganizationAddUserPage(org, userName))
}
