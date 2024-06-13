package auth

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Member struct {
	ID             string  `json:"id"`
	Status         string  `json:"status"`
	Profile        Profile `json:"profile"`
	PrivacyStatus  string  `json:"privacyStatus"`
	ActivityStatus string  `json:"activityStatus"`
	CreatedDate    string  `json:"createdDate"`
	UpdatedDate    string  `json:"updatedDate"`
}

type Profile struct {
	Nickname string `json:"nickname"`
	Slug     string `json:"slug"`
	Photo    Image  `json:"photo"`
	Cover    Image  `json:"cover"`
	Title    string `json:"title"`
}

type Image struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type Root struct {
	Member Member `json:"member"`
}

func ValidateWixUser(c *gin.Context) *Root {
	memberId := c.GetHeader("X-Wix-Member-ID")
	url := "https://www.wixapis.com/members/v1/members/" + memberId
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println("failed to create request", err)
		c.JSON(400, gin.H{"error": "Failed to create request"})
		c.Abort()
		return nil
	}
	request.Header.Set("Authorization", "Bearer "+os.Getenv("WIX_TOKEN"))
	request.Header.Set("wix-site-id", os.Getenv("WIX_SITE_ID"))
	resp, err := client.Do(request)

	if err != nil {
		println("failed to make request", err)
		c.JSON(400, gin.H{"error": "Failed to make request"})
		c.Abort()
		return nil
	}
	defer resp.Body.Close()

	var member Root
	if err := json.NewDecoder(resp.Body).Decode(&member); err != nil {
		println("failed to parse", err)
		c.JSON(400, gin.H{"error": "Failed to parse response"})
		c.Abort()
		return nil
	}

	if member.Member.ID == "" {
		println("No member found")
		c.JSON(400, gin.H{"error": "No member found"})
		c.Abort()
		return nil
	}
	return &member
}
