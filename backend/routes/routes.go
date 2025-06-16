package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ganesh96/simple-reddit/backend/posts"
	"github.com/ganesh96/simple-reddit/backend/users"
	"github.com/ganesh96/simple-reddit/backend/communities"
	// ... other imports
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Public routes
	router.POST("/login", users.Login)
	router.POST("/user", users.CreateUser)
	router.GET("/posts", posts.GetAllPosts)
	router.GET("/post/:id", posts.GetPostById)
	router.GET("/communities", communities.GetAllCommunities)
	router.GET("/community/:name", communities.GetCommunityByName)

	// Authenticated routes
	authorized := router.Group("/")
	authorized.Use(users.AuthorizeJWT())
	{
		authorized.POST("/post", posts.CreatePost)
		authorized.PATCH("/post/:id", posts.UpdatePost)
		authorized.DELETE("/post/:id", posts.DeletePost)

		authorized.POST("/community", communities.CreateCommunity)
		authorized.DELETE("/community/:name", communities.DeleteCommunity)

		authorized.GET("/profile/:username", users.GetProfile)
		authorized.DELETE("/user/:username", users.DeleteUser)
	}

	return router
}