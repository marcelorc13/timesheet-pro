package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marcelorc13/timesheet-pro/internal/domain"
	service "github.com/marcelorc13/timesheet-pro/internal/services"
	"github.com/marcelorc13/timesheet-pro/internal/templates/components"
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
		c.Status(http.StatusBadRequest)
		utils.Render(c.Request.Context(), c.Writer, components.Response(err.Error(), true))
		return
	}

	err := h.service.Create(c.Request.Context(), usuario)
	if err != nil {
		c.Status(http.StatusBadRequest)
		utils.Render(c.Request.Context(), c.Writer, components.Response(err.Error(), true))
		return
	}
	c.Header("HX-Redirect", "/login")
	c.Status(http.StatusCreated)
	utils.Render(c.Request.Context(), c.Writer, components.Response("Usuário criado com sucesso! realize o login", false))
}
func (h UserHandler) Login(c *gin.Context) {
	var usuario domain.LoginUser

	if err := c.ShouldBind(&usuario); err != nil {
		c.Status(400)
		utils.Render(c.Request.Context(), c.Writer, components.Response(err.Error(), true))
	}

	u, err := h.service.Login(c.Request.Context(), usuario)

	if err != nil {
		if strings.Contains(err.Error(), "usuário não encontrado") {
			c.Status(401)
			utils.Render(c.Request.Context(), c.Writer, components.Response("Usuário não existe", true))
			return
		}
		if strings.Contains(err.Error(), "senha incorreta") {
			c.Status(401)
			utils.Render(c.Request.Context(), c.Writer, components.Response("Senha incorreta", true))
			return
		}
		c.Status(401)
		utils.Render(c.Request.Context(), c.Writer, components.Response("Usuário ou senha incorreta", true))
		return
	}

	token, err := utils.GenerateJwtToken(u.ID.String())
	if err != nil {
		c.Status(400)
		utils.Render(c.Request.Context(), c.Writer, components.Response(fmt.Sprintf("Erro ao gerar token JWT: %v", err), true))
	}

	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.Header("HX-Redirect", "/")
	c.Status(200)
	utils.Render(c.Request.Context(), c.Writer, components.Response("Usuário logado com sucesso", false))
}
