package routes

import (
	"github.com/ganesh96/simple-reddit/backend/comments"
	"github.com/ganesh96/simple-reddit/backend/communities"
	"github.com/ganesh96/simple-reddit/backend/posts"
	"github.com/ganesh96/simple-reddit/backend/profiles"
	"github.com/ganesh96/simple-reddit/backend/users"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// User routes
	router.POST("/signup", users.Signup)
	router.POST("/login", users.Login)
	router.DELETE("/users/:username", users.AuthorizeJWT(), users.DeleteUser)

	// Profile routes
	router.GET("/profiles/:username", profiles.GetProfileByUsername)
	router.PUT("/profiles/:username", users.AuthorizeJWT(), profiles.UpdateProfile)

	// Community routes
	router.POST("/communities", users.AuthorizeJWT(), communities.CreateCommunity)
	router.GET("/communities", communities.GetAllCommunities)
	router.DELETE("/communities/:communityName", users.AuthorizeJWT(), communities.DeleteCommunityByName)

	// Post routes
	router.POST("/posts", users.AuthorizeJWT(), posts.CreatePost)
	router.GET("/posts", posts.GetAllPosts)
	router.GET("/posts/:postId", posts.GetPostById)
	router.PUT("/posts/:postId", users.AuthorizeJWT(), posts.UpdatePost)
	router.DELETE("/posts/:postId", users.AuthorizeJWT(), posts.DeletePost)

	// Comment routes
	router.POST("/posts/:postId/comments", users.AuthorizeJWT(), comments.CreateComment)
	router.GET("/posts/:postId/comments", comments.GetCommentsByPostId)
	router.PUT("/comments/:commentId", users.AuthorizeJWT(), comments.UpdateComment)
	router.DELETE("/comments/:commentId", users.AuthorizeJWT(), comments.DeleteComment)
}
