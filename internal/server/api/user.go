package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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

	token, err := utils.GenerateJwtToken(u.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.HttpResponse{Status: http.StatusInternalServerError, Message: fmt.Sprintf("Erro ao gerar token JWT: %v", err)})
		return
	}

	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.Header("HX-Redirect", "/")
	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Login bem-sucedido"})
}
