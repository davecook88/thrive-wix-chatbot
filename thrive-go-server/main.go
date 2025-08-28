package main

import (
	"context"
	"fmt"
	"net/http"
	"thrive/server/admin"
	chromeext "thrive/server/chrome-ext"
	"thrive/server/db"
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
		c.Writer.Header().Set("Access-control-allow-headers", "Content-Type, Authorization")

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

	dbClient, err := db.NewClient(context.Background())
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

	r.POST("/chat", handlers.PostMessageHandler(dbClient))
	r.GET("/chat", handlers.GetChatHandler(dbClient))

	admin.RegisterAdminRoutes(r, dbClient)
	chromeext.RegisterChromeExtRoutes(r)

	r.Run() // listen and serve on 0.0.0.0:8080
}
