package main

import (
	"fmt"
	"net/http"
	"thrive/server/admin"
	chromeext "thrive/server/chrome-ext"
	"thrive/server/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		fmt.Println("origin", origin)
		// if origin starts with os.Getenv("ALLOWED_ORIGIN"), set the header
		// if origin != "" && strings.HasPrefix(origin, os.Getenv("ALLOWED_ORIGIN")) {
		// 	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		// }
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.Use(corsMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/chat", handlers.PostMessageHandler)
	r.GET("/chat", handlers.GetChatHandler)

	admin.RegisterAdminRoutes(r)
	chromeext.RegisterChromeExtRoutes(r)

	r.Run() // listen and serve on 0.0.0.0:8080
}
