package routes

import (
	"github.com/ganesh96/simple-reddit/backend/comments"
	"github.com/ganesh96/simple-reddit/backend/communities"
	"github.com/ganesh96/simple-reddit/backend/posts"
	"github.com/ganesh96/simple-reddit/backend/profiles"
	"github.com/ganesh96/simple-reddit/backend/users"
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {
	// Public routes
	router.POST("/login", users.Login)
	router.POST("/user", users.CreateUser)
	router.GET("/posts", posts.GetAllPosts)
	router.GET("/post/:id", posts.GetPostById)
	router.GET("/communities", communities.GetAllCommunities)
	router.GET("/community/:name", communities.GetCommunityByName)
	router.GET("/comments/post/:postid", comments.GetCommentsForPost())

	// Authenticated routes
	authorized := router.Group("/")
	authorized.Use(users.AuthorizeJWT())
	{
		authorized.POST("/post", posts.CreatePost)
		authorized.PATCH("/post/:id", posts.UpdatePost)
		authorized.DELETE("/post/:id", posts.DeletePost)

		authorized.POST("/community", communities.CreateCommunity)
		authorized.DELETE("/community/:name", communities.DeleteCommunity)

		authorized.GET("/profile/:username", profiles.GetProfile())
		authorized.DELETE("/user/:username", users.DeleteUser)

		authorized.POST("/comment", comments.CreateComment())
		authorized.DELETE("/comment/:id", comments.DeleteComment())
	}
}
