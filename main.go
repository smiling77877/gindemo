package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	c := &UserHandler{}
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
