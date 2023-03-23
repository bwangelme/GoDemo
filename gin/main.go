package main

import (
	"bwdemo/wtypes"
	"fmt"
	"github.com/gin-gonic/gin"
)

func BindDemo(c *gin.Context) {
	var args struct {
		IsTv wtypes.Bool `json:"is_tv"`
	}
	err := c.BindJSON(&args)
	fmt.Println(args.IsTv)
	c.JSON(200, gin.H{
		"args": args,
		"err":  err,
	})
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {
	gin.SetMode("debug")
	r := gin.Default()
	r.GET("/ping", Ping)
	r.POST("/bind", BindDemo)
	r.Run()
}
