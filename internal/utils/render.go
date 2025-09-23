package utils

import (
	"context"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func Render(ctx context.Context, w gin.ResponseWriter, comp templ.Component) error {
	return comp.Render(ctx, w)
}
