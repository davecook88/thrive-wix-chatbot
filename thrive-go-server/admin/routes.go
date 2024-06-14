package admin

import (
	"strconv"
	"thrive/server/auth"
	"thrive/server/db"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r *gin.Engine) {
	adminGroup := r.Group("/admin")
	adminGroup.Use(auth.ValidateWixHeader())
	{
		adminGroup.GET("/list-chats", listChats)
	}
}

func listChats(c *gin.Context) {
	dbClient, err := db.NewClient(c, "thrive-chat")
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create db client"})
		return
	}
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 1000
	}
	offset, _ := strconv.Atoi(c.Query("offset"))
	chats := dbClient.ListChats(c, &db.ListChatsParams{
		Limit:  limit,
		Offset: offset,
	})
	c.JSON(200, gin.H{"chats": chats})

}
