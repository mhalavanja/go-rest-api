package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()

	authorized := router.Group("")
	authorized.Use(gin.BasicAuth(gin.Accounts{
		"admin": "admin",
	}))

	authorized.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run()
}
