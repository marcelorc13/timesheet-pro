// Package api provides HTTP handlers for REST API endpoints.
// It handles request parsing, validation, and response formatting for the application.
package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/marcelorc13/timesheet-pro/internal/domain"
	service "github.com/marcelorc13/timesheet-pro/internal/services"
	"github.com/marcelorc13/timesheet-pro/internal/utils"
)

type OrganizationHandler struct {
	service service.OrganizationService
}

func NewOrganizationHandler(us service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{us}
}

func (h OrganizationHandler) Create(c *gin.Context) {
	var co domain.CreateOrganization

	if err := c.ShouldBind(&co); err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	orgID, err := h.service.CreateWithUser(c.Request.Context(), co)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}
	c.Header("HX-Redirect", "/")
	c.JSON(http.StatusCreated, domain.HttpResponse{Status: http.StatusCreated, Message: "Organização criada com sucesso", Data: orgID})
}

func (h OrganizationHandler) List(c *gin.Context) {
	// Extract user ID from JWT token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token não fornecido"})
		return
	}

	claims, err := utils.GetTokenClaims(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Token inválido"})
		return
	}

	userIDStr, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "ID de usuário inválido no token"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID de usuário inválido"})
		return
	}

	// Get user's organization
	res, err := h.service.GetOrganizationByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.HttpResponse{Status: http.StatusNotFound, Message: "Usuário não pertence a nenhuma organização"})
		return
	}

	if res == nil {
		c.JSON(http.StatusNotFound, domain.HttpResponse{Status: http.StatusNotFound, Message: "Organização não encontrada"})
		return
	}

	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Organização do usuário", Data: res})
}

func (h OrganizationHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	orgID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID inválido"})
		return
	}

	res, err := h.service.GetByID(c.Request.Context(), orgID)

	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	if res == nil {
		c.JSON(http.StatusNotFound, domain.HttpResponse{Status: http.StatusNotFound, Message: "Organização não encontrada"})
		return
	}

	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: fmt.Sprintf("Organização de id %s encontrada", id), Data: res})
}

func (h OrganizationHandler) GetByUserID(c *gin.Context) {
	userId := c.Param("userId")

	userID, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID de usuário inválido"})
		return
	}

	res, err := h.service.GetOrganizationByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.HttpResponse{Status: http.StatusNotFound, Message: "Usuário não pertence a nenhuma organização"})
		return
	}

	if res == nil {
		c.JSON(http.StatusNotFound, domain.HttpResponse{Status: http.StatusNotFound, Message: "Organização não encontrada"})
		return
	}

	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Organização do usuário", Data: res})
}

func (h OrganizationHandler) Update(c *gin.Context) {
	id := c.Param("id")

	orgID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID inválido"})
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

	var uo domain.UpdateOrganization
	if err := c.ShouldBind(&uo); err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	err = h.service.Update(c.Request.Context(), userID, orgID, uo)
	if err != nil {
		if err.Error() == "usuário não tem permissão para atualizar esta organização" {
			c.JSON(http.StatusForbidden, domain.HttpResponse{Status: http.StatusForbidden, Message: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	c.Header("HX-Redirect", "/")
	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Organização atualizada com sucesso"})
}

func (h OrganizationHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	orgID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID inválido"})
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

	err = h.service.Delete(c.Request.Context(), userID, orgID)
	if err != nil {
		if err.Error() == "usuário não tem permissão para deletar esta organização" {
			c.JSON(http.StatusForbidden, domain.HttpResponse{Status: http.StatusForbidden, Message: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Organização deletada com sucesso"})
}

func (h OrganizationHandler) AddUser(c *gin.Context) {
	id := c.Param("id")

	orgID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID inválido"})
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

	var addUser domain.AddUserToOrganization
	if err := c.ShouldBind(&addUser); err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	err = h.service.AddUserByEmail(c.Request.Context(), userID, orgID, addUser)
	if err != nil {
		if err.Error() == "usuário não tem permissão para adicionar usuários a esta organização" {
			c.JSON(http.StatusForbidden, domain.HttpResponse{Status: http.StatusForbidden, Message: err.Error()})
			return
		}
		if err.Error() == "usuário já é membro desta organização" {
			c.JSON(http.StatusConflict, domain.HttpResponse{Status: http.StatusConflict, Message: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	c.Header("HX-Redirect", "/")
	c.JSON(http.StatusCreated, domain.HttpResponse{Status: http.StatusCreated, Message: "Usuário adicionado à organização com sucesso"})
}

func (h OrganizationHandler) RemoveUser(c *gin.Context) {
	id := c.Param("id")

	orgID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID da organização inválido"})
		return
	}

	utID := c.Param("userId")
	userTargetID, err := uuid.Parse(utID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID do usuário inválido"})
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

	err = h.service.RemoveUserFromOrganization(c.Request.Context(), userID, orgID, userTargetID)
	if err != nil {
		if err.Error() == "usuário não tem permissão para remover membros desta organização" {
			c.JSON(http.StatusForbidden, domain.HttpResponse{Status: http.StatusForbidden, Message: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	c.Header("HX-Refresh", "true")
	c.JSON(http.StatusNoContent, domain.HttpResponse{Status: http.StatusNoContent, Message: "Usuário removido da organização com sucesso"})
}

func (h OrganizationHandler) Leave(c *gin.Context) {
	id := c.Param("id")

	orgID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID inválido"})
		return
	}

	// Extract user ID from JWT token
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
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "ID de usuário inválido no token"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID de usuário inválido"})
		return
	}

	err = h.service.LeaveOrganization(c.Request.Context(), userID, orgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	c.Header("HX-Redirect", "/")
	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Você saiu da organização com sucesso"})
}
