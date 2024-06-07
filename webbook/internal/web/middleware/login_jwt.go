package middleware

import (
	webook "gindemo/webbook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type LoginJWTMiddlewareBuilder struct {
	webook.Handler
}

func NewLoginJWTMiddlewareBuilder(hdl webook.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: hdl,
	}
}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" ||
			path == "/users/login" ||
			path == "/users/login_sms/code/send" ||
			path == "/users/login_sms" ||
			path == "/oauth2/wechat/authurl" ||
			path == "/oauth2/wechat/callback" {
			// 不需要登录校验
			return
		}

		tokenStr := m.ExtractToken(ctx)
		var uc webook.UserClaims

		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return webook.JWTKey, nil
		})
		if err != nil {
			// token不对
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			// 在这里发现 access_token 过期了, 生成一个新的 access_token

			// token解析出来了，但是token可能是非法的，或者过期了的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 这里看
		err = m.CheckSession(ctx, uc.Ssid)
		if err != nil {
			// token无效或者 redis 有问题
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 可以兼容redis异常的情况
		// 做好监控,监控有没有error

		ctx.Set("user", uc)
	}
}
