package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"nsq-demoset/app/_applib/appresponse"
	"nsq-demoset/app/_applib/utils"
	"nsq-demoset/app/app-services/cmd/front_api/middleware"
	"nsq-demoset/app/app-services/conf"
	"nsq-demoset/app/app-services/model"
	"nsq-demoset/app/nsq-services/ds"
)

type authHandler struct {
	R        *gin.Engine
	UserSvc  model.UserService
	TokenSvc model.TokenService
}

func NewAuthHandler(h *Handler) *authHandler {
	return &authHandler{
		R:        h.R,
		UserSvc:  h.userScv,
		TokenSvc: h.tokenScv,
	}
}

type reqPostLogin struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type reqRefreshToken struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
}

func (ctr *authHandler) Register() {
	group := ctr.R.Group("/api/auth")

	// guest
	group.POST("/login", ctr.postLogin)
	group.POST("/register", ctr.postRegister)

	// auth
	group.Use(middleware.AuthMiddleware(ctr.UserSvc))
	group.POST("/logout", ctr.postLogout)
	group.POST("/refresh", ctr.postRefresh)
	group.POST("/me", ctr.getMe)
}

// @BasePath /api/auth

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description User Login
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Login
// @Router /login [post]
func (ctr *authHandler) postLogin(c *gin.Context) {
	res := &appresponse.ResponseObj{}
	req := &reqPostLogin{}

	if err := c.ShouldBind(&req); err != nil {
		res.ErrCode = 422
		res.ErrMsg = err.Error()
		c.JSON(422, res)
		return
	}

	ctx := c.Request.Context()
	user, err := ctr.UserSvc.FindByEmail(ctx, req.Email)
	if err != nil {
		res.ErrCode = 404
		res.ErrMsg = "User is not register yet"
		c.JSON(404, res)
		return
	}

	//logger.Sugar.Debug(user)
	//
	//if !user.Activated {
	//	res.ErrCode = 401
	//	res.ErrMsg = "Please activate first to login"
	//	c.JSON(401, res)
	//	return
	//}

	// validate password
	ok, err := utils.ComparePasswords(user.Password, req.Password)
	if err != nil {
		res.ErrCode = 500
		res.ErrMsg = err.Error()
		c.JSON(500, res)
		return
	}

	if !ok {
		res.ErrCode = 401
		res.ErrMsg = "Wrong Password"
		c.JSON(401, res)
		return
	}

	// Generate token pair
	tokenPair, err := ctr.TokenSvc.GenerateTokenPair(ctx, user, "")
	if err != nil {
		res.ErrCode = 0
		res.ErrMsg = "Error on generating token pair"
		c.JSON(200, res)
		return
	}

	res.ErrCode = 0
	res.ErrMsg = "Success"
	res.Data = gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	}
	c.JSON(200, res)
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description User Register
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Register
// @Router /register [post]
func (ctr *authHandler) postRegister(c *gin.Context) {
	res := &appresponse.ResponseObj{}
	req := &reqPostLogin{}

	if err := c.ShouldBind(&req); err != nil {
		res.ErrCode = http.StatusUnprocessableEntity
		res.ErrMsg = err.Error()
		c.JSON(422, res)
		return
	}

	// hash password
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		res.ErrCode = 500
		res.ErrMsg = err.Error()
		c.JSON(500, res)
		return
	}

	ctx := c.Request.Context()

	// create user
	user := &model.User{}
	user.Email = req.Email
	user.Password = hashPassword
	_, err = ctr.UserSvc.Create(ctx, user)
	if err != nil {
		res.ErrCode = 500
		res.ErrMsg = err.Error()
		c.JSON(500, res)
		return
	}

	res.ErrCode = 0
	res.ErrMsg = "Success"
	c.JSON(200, res)
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description User Logout
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Logout
// @Router /logout [post]
func (ctr *authHandler) postLogout(c *gin.Context) {
	res := &appresponse.ResponseObj{}
	user := c.MustGet("user").(*model.User)

	req := &reqRefreshToken{}
	if err := c.ShouldBind(&req); err != nil {
		res.ErrCode = 422
		res.ErrMsg = err.Error()
		c.JSON(422, res)
		return
	}

	refreshClaims, err := utils.ValidateRefreshToken(req.RefreshToken, conf.RefreshSecret)
	if err != nil {
		res.ErrCode = 401
		res.ErrMsg = err.Error()
		c.JSON(401, res)
		return
	}

	ctx := c.Request.Context()
	key := fmt.Sprintf("%s:%v", refreshClaims.Id, user.Id)

	// remove old refresh token from redis
	_, err = ds.RDB.Del(ctx, key).Result()
	if err != nil {
		res.ErrCode = 500
		res.ErrMsg = err.Error()
		c.JSON(500, res)
		return
	}

	res.ErrCode = 0
	res.ErrMsg = "success"
	c.JSON(200, res)
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description User Refresh
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Refresh
// @Router /refresh [post]
func (ctr *authHandler) postRefresh(c *gin.Context) {
	res := &appresponse.ResponseObj{}
	user := c.MustGet("user").(*model.User)
	ctx := c.Request.Context()

	req := reqRefreshToken{}
	if err := c.ShouldBind(&req); err != nil {
		res.ErrCode = 401
		res.ErrMsg = err.Error()
		c.JSON(401, res)
		return
	}

	_, err := utils.ValidateRefreshToken(req.RefreshToken, conf.RefreshSecret)
	if err != nil {
		res.ErrCode = 401
		res.ErrMsg = "Refresh token is expired"
		c.JSON(401, res)
		return
	}

	// generate token pair
	tokenPair, err := ctr.TokenSvc.GenerateTokenPair(ctx, user, req.RefreshToken)
	if err != nil {
		res.ErrCode = 500
		res.ErrMsg = "Error on generating token pair"
		c.JSON(500, res)
		return
	}

	res.ErrCode = 0
	res.ErrMsg = "Success"
	res.Data = gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	}
	c.JSON(http.StatusOK, res)
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description User Me
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Me
// @Router /me [post]
func (ctr *authHandler) getMe(c *gin.Context) {
	res := &appresponse.ResponseObj{}
	user := c.MustGet("user").(*model.User)

	res.ErrCode = 0
	res.ErrMsg = "Success"
	res.Data = user
	c.JSON(http.StatusOK, res)
}
