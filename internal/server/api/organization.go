package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marcelorc13/timesheet-pro/internal/domain"
	service "github.com/marcelorc13/timesheet-pro/internal/services"
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
	c.Header("HX-Redirect", "/login")
	c.JSON(http.StatusCreated, domain.HttpResponse{Status: http.StatusCreated, Message: "Organização criada com sucesso", Data: orgID})
}

