package configs

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const devJWTSecret = "dev-only-change-me"

func signingKey() []byte {
	secret := SecretKey()
	if secret == "" {
		secret = devJWTSecret
	}
	return []byte(secret)
}

// JWTClaim adds a username to the standard claims.
type JWTClaim struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a given username.
func GenerateToken(username string) (string, error) {
	claims := &JWTClaim{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey())
}

// ValidateToken validates the jwt token.
func ValidateToken(signedToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return signingKey(), nil
		},
	)
}
