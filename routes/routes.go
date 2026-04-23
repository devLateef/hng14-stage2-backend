package routes

import (
	"insight-api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	api := r.Group("/api")
	{
		api.GET("/profiles", handlers.GetProfiles)
		api.GET("/profiles/search", handlers.SearchProfiles)
	}
}
