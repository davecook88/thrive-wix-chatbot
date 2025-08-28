package admin

import (
	"strconv"
	"thrive/server/auth"
	"thrive/server/db"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r *gin.Engine, dbClient *db.Client) {
	adminGroup := r.Group("/admin")
	adminGroup.Use(auth.ValidateWixHeader())
	{
		adminGroup.GET("/list-chats", listChats(dbClient))
	}
}

func listChats(dbClient *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
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

}
