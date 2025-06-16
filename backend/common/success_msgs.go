package common

const (
	// SUCCESS is a generic success code
	SUCCESS = "SUCCESS"
	// CREATED is for successful resource creation
	CREATED = "CREATED"
)

// SuccessMessages maps success codes to their messages
var SuccessMessages = map[string]APIMessage{
	SUCCESS: {Message: "Operation successful", Code: SUCCESS},
	CREATED: {Message: "Resource created successfully", Code: CREATED},
}

// ErrorMessages maps error codes to their messages
var ErrorMessages = map[string]APIMessage{
	MONGO_DB_ERROR:             {Message: "Database error", Code: MONGO_DB_ERROR},
	INVALID_REQUEST_BODY:       {Message: "Invalid request body", Code: INVALID_REQUEST_BODY},
	INVALID_PARAM:              {Message: "Invalid URL parameter", Code: INVALID_PARAM},
	UNAUTHORIZED:               {Message: "Unauthorized access", Code: UNAUTHORIZED},
	EMAIL_ALREADY_EXISTS:       {Message: "Email already exists", Code: EMAIL_ALREADY_EXISTS},
	USERNAME_ALREADY_EXISTS:   {Message: "Username already exists", Code: USERNAME_ALREADY_EXISTS},
	USER_NOT_FOUND:             {Message: "User not found", Code: USER_NOT_FOUND},
	INCORRECT_PASSWORD:         {Message: "Incorrect password", Code: INCORRECT_PASSWORD},
	COMMUNITY_NOT_FOUND:        {Message: "Community not found", Code: COMMUNITY_NOT_FOUND},
	COMMUNITY_ALREADY_EXISTS: {Message: "Community with this name already exists", Code: COMMUNITY_ALREADY_EXISTS},
	POST_NOT_FOUND:             {Message: "Post not found", Code: POST_NOT_FOUND},
	COMMENT_NOT_FOUND:          {Message: "Comment not found", Code: COMMENT_NOT_FOUND},
}
