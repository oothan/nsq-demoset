package handler

import (
	"github.com/gin-gonic/gin"
	"nsq-demoset/app/app-services/cmd/front_api/middleware"
	"nsq-demoset/app/app-services/ds"
	"nsq-demoset/app/app-services/model"
	"nsq-demoset/app/app-services/repository"
	"nsq-demoset/app/app-services/service"
)

type Handler struct {
	R *gin.Engine

	userScv  model.UserService
	tokenScv model.TokenService
}

type HConfig struct {
	R  *gin.Engine
	DS *ds.DataSource
}

func NewHandler(c *HConfig) *Handler {

	// token
	tokenRepo := repository.NewTokenRepository(c.DS)
	tokenService := service.NewTokenService(&service.TokenConfig{
		TokenRepo: tokenRepo,
	})

	// user repo
	userRepo := repository.NewUserRepository(c.DS)
	userService := service.NewUserService(&service.UserConfig{
		UserRepo: userRepo,
	})

	return &Handler{
		R:        c.R,
		userScv:  userService,
		tokenScv: tokenService,
	}
}

func (h *Handler) Register() {
	// register cors middleware
	h.R.Use(middleware.Cors())

	// home
	homeHandler := NewHomeHandler(h)
	homeHandler.Register()

	// auth
	authHandler := NewAuthHandler(h)
	authHandler.Register()

}
