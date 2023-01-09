package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"nsq-demoset/app/_applib"
	logger "nsq-demoset/app/_applib"
	"nsq-demoset/app/app-services/cmd/front_api/middleware"
	"nsq-demoset/app/app-services/internal/dto"
	"nsq-demoset/app/app-services/internal/model"
)

type postHandler struct {
	R       *gin.Engine
	PostSvc model.PostService
	UserSvc model.UserService
}

func NewPostHandler(h *Handler) *postHandler {
	return &postHandler{
		R:       h.R,
		PostSvc: h.postScv,
		UserSvc: h.userScv,
	}
}

func (ctr *postHandler) Register() {
	group := ctr.R.Group("/api/post")

	// auth
	group.Use(middleware.AuthMiddleware(ctr.UserSvc))
	group.POST("/view", ctr.postView)
}

func (ctr *postHandler) postView(c *gin.Context) {
	res := &dto.ResponseObj{}
	user := c.MustGet("user").(*model.User)
	logger.Sugar.Debug("UserId =========== ", user.Id)

	req := dto.RequestPostView{}
	if err := c.ShouldBind(&req); err != nil {
		res.ErrCode = _applib.ErrorCode401
		res.ErrMsg = err.Error()
		c.JSON(_applib.ErrorCode401, res)
		return
	}

	post, err := ctr.PostSvc.FindById(context.Background(), req.PostId)
	if err != nil {
		res.ErrCode = _applib.ErrorCode505
		res.ErrMsg = err.Error()
		c.JSON(_applib.ErrorCode401, res)
		return
	}

	bytes, err := json.Marshal(post)

	postRes := dto.PostResponse{}
	if err := json.Unmarshal([]byte(bytes), &postRes); err != nil {
		res.ErrCode = _applib.ErrorCode505
		res.ErrMsg = err.Error()
		c.JSON(_applib.ErrorCode401, res)
		return
	}

	res.ErrCode = _applib.SuccessCode
	res.ErrMsg = _applib.SuccessMessage
	res.Data = postRes
	c.JSON(http.StatusOK, res)
}
