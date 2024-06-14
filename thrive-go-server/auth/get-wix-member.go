package auth

import (
	"thrive/server/wix"

	"github.com/gin-gonic/gin"
)

func ValidateWixUser(c *gin.Context) *wix.WixMember {
	memberId := c.GetHeader("X-Wix-Member-ID")

	client := wix.NewWixClient()
	member, err := client.GetMember(memberId)

	if err != nil {
		println("failed to make request", err)
		c.JSON(400, gin.H{"error": "Failed to make request"})
		c.Abort()
		return nil
	}

	if member.ID == "" {
		println("No member found")
		c.JSON(400, gin.H{"error": "No member found"})
		c.Abort()
		return nil
	}
	return member
}
