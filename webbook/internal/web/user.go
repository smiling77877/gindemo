package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
}

func (c *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/users/signup", c.Signup)
	//server.POST("/users/login", c.Login)
	//server.POST("/users/edit", c.Edit)
	//server.GET("/users/profile", c.Profile)
}

func (c *UserHandler) Signup(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello, 你在注册")
}
