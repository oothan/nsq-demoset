package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type homeHandler struct {
	R *gin.Engine
}

func NewHomeHandler(h *Handler) *homeHandler {
	return &homeHandler{
		R: h.R,
	}
}

// Register PingExample godoc
// @Summary ping example
// @Schemes
// @Description Home
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Welcome
// @router / [get]
func (ctr *homeHandler) Register() {
	ctr.R.GET("/", ctr.welcome)
}

func (ctr *homeHandler) welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome Home ",
	})
}
