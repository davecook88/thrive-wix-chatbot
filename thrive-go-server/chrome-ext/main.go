package chromeext

import (
	"github.com/gin-gonic/gin"
)

func RegisterChromeExtRoutes(r *gin.Engine) {
	adminGroup := r.Group("/ext")
	adminGroup.Use(ValidateChromeExtHeader())
	{
		adminGroup.GET("/test", test)
	}
}

func test(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
