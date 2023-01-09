package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"nsq-demoset/app/app-services/conf"
	"nsq-demoset/app/app-services/internal/dto"
	"nsq-demoset/app/app-services/internal/model"
	"nsq-demoset/app/app-services/internal/utils"
	"strings"
	"time"
)

type authHeader struct {
	AccessToken string `header:"Authorization"`
}

func AuthMiddleware(userSvc model.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}
		res := &dto.ResponseObj{}

		// bind Authorization Header to h and check for validation errors
		if err := c.ShouldBindHeader(&h); err != nil {
			res.ErrCode = 401
			res.ErrMsg = "Must provide `Authorization` header in format of `Bearer {token}`"
			c.JSON(401, res)
			c.Abort()
			return
		}

		accessToken := strings.Split(h.AccessToken, "Bearer ")

		if len(accessToken) < 2 {
			res.ErrCode = 401
			res.ErrMsg = "Must provide `Authorization` in format of `Bearer {token}`"
			c.JSON(401, res)
			c.Abort()
			return
		}

		// validate access token
		accessTokenClaim, err := utils.ValidateAccessToken(accessToken[1], conf.PublicKey)
		if err != nil {
			res.ErrCode = 403
			res.ErrMsg = "Permission Denied"
			c.JSON(403, res)
			c.Abort()
			return
		}

		uid := accessTokenClaim.User.Id
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		user, err := userSvc.FindById(ctx, uid)
		if err != nil {
			res.ErrCode = 403
			res.ErrMsg = "Permission Denied"
			c.JSON(403, res)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
