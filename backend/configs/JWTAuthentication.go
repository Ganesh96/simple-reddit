package configs

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// mySigningKey is used to sign the JWT token.
// It's better to load this from environment variables for production.
var mySigningKey = []byte("supersecretkey")

// JWTClaim adds a username to the standard claims
type JWTClaim struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a given username
func GenerateToken(username string) (string, error) {
	// Create the claims
	claims := &JWTClaim{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expires after 24 hours
		},
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates the jwt token
func ValidateToken(signedToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return mySigningKey, nil
		},
	)
}
