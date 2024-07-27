package chromeext

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var CHROME_EXT_TOKEN_HEADER = "X-Chrome-Ext-Token"

func ValidateChromeExtHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		instance := c.GetHeader(CHROME_EXT_TOKEN_HEADER)
		if instance == "" {
			c.JSON(400, gin.H{"error": fmt.Sprintf("Missing %s header", CHROME_EXT_TOKEN_HEADER)})
			c.Abort()
			return
		}
		appSecret := os.Getenv("CHROME_EXT_TOKEN")
		if appSecret == "" {
			c.JSON(500, gin.H{"error": "No app secret"})
			c.Abort()
			return
		}

		c.Next()
	}
}
