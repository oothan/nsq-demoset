package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nsq-demoset/app/_applib"
	"nsq-demoset/app/_applib/nsq"
	"nsq-demoset/app/app-services/internal/dto"
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
// @Tags Home
// @Accept json
// @Produce json
// @router / [get]
func (ctr *homeHandler) Register() {
	ctr.R.POST("/", ctr.welcome)
}

func (ctr *homeHandler) welcome(c *gin.Context) {
	res := &dto.ResponseObj{}
	req := dto.HomeRequest{}
	if err := c.ShouldBind(&req); err != nil {
		res.ErrCode = _applib.ErrorCode401
		res.ErrMsg = err.Error()
		c.JSON(_applib.ErrorCode401, res)
		return
	}

	nsq.NsqTestMessageEvent(req.Message)

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome Home ",
	})
}
