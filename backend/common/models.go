// File: backend/common/models.go
// This is a new file to hold shared structs.
package common

// APIMessage defines the structure for a standard API message string.
type APIMessage string

// File: backend/common/error_msgs.go
package common

const (
	// Error messages
	INVALID_REQUEST_BODY     APIMessage = "Invalid request body"
	MONGO_DB_ERROR           APIMessage = "Database error"
	INVALID_POST_ID          APIMessage = "Invalid post ID"
	POST_NOT_FOUND           APIMessage = "Post not found"
	COMMUNITY_ALREADY_EXISTS APIMessage = "Community with that name already exists"
	COMMUNITY_NOT_FOUND      APIMessage = "Community not found"
	USERNAME_ALREADY_EXISTS  APIMessage = "Username already exists"
	INVALID_CREDENTIALS      APIMessage = "Invalid username or password"
	USER_NOT_FOUND           APIMessage = "User not found"
	FORBIDDEN                APIMessage = "Forbidden"
)
