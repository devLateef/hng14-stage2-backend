package routes

import (
	"net/http"

	"insight-api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// CORS middleware — must run before any route matching
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	api := r.Group("/api")
	{
		// NOTE: /profiles/search must be registered BEFORE /profiles
		// so Gin doesn't treat "search" as a dynamic segment.
		api.GET("/profiles/search", handlers.SearchProfiles)
		api.GET("/profiles", handlers.GetProfiles)
	}
}
