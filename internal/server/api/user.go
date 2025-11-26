package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/marcelorc13/timesheet-pro/internal/domain"
	service "github.com/marcelorc13/timesheet-pro/internal/services"
	"github.com/marcelorc13/timesheet-pro/internal/utils"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{us}
}

func (h UserHandler) List(c *gin.Context) {
	res, err := h.service.List(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if res == nil {
		c.JSON(http.StatusNotFound, domain.HttpResponse{Status: http.StatusNotFound, Message: "O banco ainda não possui usuários"})
		return
	}
	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Todos os usuarios do banco", Data: res})
}

func (h UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	res, err := h.service.GetByID(c.Request.Context(), id)

	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	if res == nil {
		c.JSON(http.StatusNotFound, domain.HttpResponse{Status: http.StatusNotFound, Message: "Usuário não encontrado"})
		return
	}

	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: fmt.Sprintf("Usuário de id %s encontrado", id), Data: res})
}

func (h UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	convertId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(400, domain.HttpResponse{Status: http.StatusBadGateway, Message: "Erro ao converter número de id para inteiro"})
	}

	err = h.service.Delete(c.Request.Context(), convertId)

	if err != nil {
		if strings.Contains(err.Error(), "usuário não encontrado") {
			c.JSON(http.StatusNotFound, domain.HttpResponse{Status: http.StatusNotFound, Message: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Usuário deletado com sucesso"})
}

func (h UserHandler) Create(c *gin.Context) {
	var usuario domain.User

	if err := c.ShouldBind(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	err := h.service.Create(c.Request.Context(), usuario)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}
	c.Header("HX-Redirect", "/login")
	c.JSON(http.StatusCreated, domain.HttpResponse{Status: http.StatusCreated, Message: "Usuário criado com sucesso"})
}
func (h UserHandler) Login(c *gin.Context) {
	var usuario domain.LoginUser

	if err := c.ShouldBind(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	u, err := h.service.Login(c.Request.Context(), usuario)

	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Usuário ou senha incorreta"})
		return
	}

	token, err := utils.GenerateJwtToken(u.ID.String(), u.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.HttpResponse{Status: http.StatusInternalServerError, Message: fmt.Sprintf("Erro ao gerar token JWT: %v", err)})
		return
	}

	// Set cookie with production-ready settings
	// Domain is empty ("") to work on any domain (localhost or render.com)
	// Secure is true for HTTPS in production (Render automatically uses HTTPS)
	// HttpOnly is true to prevent XSS attacks
	c.SetCookie("token", token, 3600, "/", "", true, true)
	c.Header("HX-Redirect", "/")
	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Login bem-sucedido"})
}

// GetMyProfile gets the current user's profile from JWT token
func (h UserHandler) GetMyProfile(c *gin.Context) {
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

	user, err := h.service.GetByID(c.Request.Context(), userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, domain.HttpResponse{Status: http.StatusNotFound, Message: "Usuário não encontrado"})
		return
	}

	// Don't send password in response
	user.Password = ""
	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Data: user})
}

// UpdateMyProfile updates the current user's profile
func (h UserHandler) UpdateMyProfile(c *gin.Context) {
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

	var req struct {
		Name  string `json:"name" form:"name" binding:"required"`
		Email string `json:"email" form:"email" binding:"required,email"`
	}


	if err = c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "Dados inválidos: " + err.Error()})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: "ID inválido"})
		return
	}

	updatedUser, err := h.service.UpdateProfile(c.Request.Context(), userID, req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	// Don't send password
	updatedUser.Password = ""
	c.Header("HX-Redirect", "/")
	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Perfil atualizado com sucesso", Data: updatedUser})
}
