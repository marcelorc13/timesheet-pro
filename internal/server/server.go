// Package server represents the Controller layer, that have direct the access to the api
package server

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Router *gin.Engine
}

func NewRouter(r *gin.Engine) *Router {
	return &Router{Router: r}
}

func (r Router) Start() {
	port := os.Getenv("PORT")

	http.ListenAndServe(port, r.Router)
}
