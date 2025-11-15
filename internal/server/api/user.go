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

func (h UserHandler) GetUsuarios(c *gin.Context) {
	res, err := h.service.GetUsuarios(c.Request.Context())

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

func (h UserHandler) GetUsuario(c *gin.Context) {
	id := c.Param("id")

	res, err := h.service.GetUsuario(c.Request.Context(), id)

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

func (h UserHandler) DeleteUsuario(c *gin.Context) {
	id := c.Param("id")

	convertId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(400, domain.HttpResponse{Status: http.StatusBadGateway, Message: "Erro ao converter número de id para inteiro"})
	}

	err = h.service.DeleteUsuario(c.Request.Context(), convertId)

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

func (h UserHandler) CreateUsuario(c *gin.Context) {
	var usuario domain.Usuario

	err := c.BindJSON(&usuario)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	err = h.service.CreateUsuario(c.Request.Context(), usuario)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain.HttpResponse{Status: http.StatusCreated, Message: "Usuário criado com sucesso"})
}
func (h UserHandler) Login(c *gin.Context) {
	var usuario domain.LoginUsuario

	if err := c.BindJSON(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status: http.StatusBadRequest, Message: err.Error()})
	}

	u, err := h.service.Login(c.Request.Context(), usuario)

	if err != nil {
		if strings.Contains(err.Error(), "usuário não encontrado") {
			c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Usuário não existe"})
			return
		}
		if strings.Contains(err.Error(), "senha incorreta") {
			c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Senha incorreta"})
			return
		}
		c.JSON(http.StatusUnauthorized, domain.HttpResponse{Status: http.StatusUnauthorized, Message: "Usuário ou senha incorreta"})
		return
	}

	token, err := utils.GenerateJwtToken(u.ID.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.HttpResponse{Status:http.StatusBadRequest, Message: fmt.Sprintf("Erro ao gerar token JWT: %v", err)})
	}

	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, domain.HttpResponse{Status: http.StatusOK, Message: "Usuário logado com sucesso"})
}
