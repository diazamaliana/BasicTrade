package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(getSecretKey()) // Replace with a strong secret key

const (
	errorSignInToProceed    = "Sign in to proceed"
	errorExpireClaimMissing = "Expire claim is missing"
	errorExpireClaimType    = "Expire claim is not a valid type"
	errorParsingExpTime     = "Error parsing expiration time"
	errorTokenExpired       = "Token is expired"
)

// getSecretKey reads the secret key from the environment variable
func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		fmt.Println("Warning: JWT_SECRET_KEY not set. Using a default key. Please set a strong secret key in production.")
		return "your-256-bit-secret" // Replace with a strong default key
	}
	return secretKey
}

func GenerateToken(adminUUID string, email string) (string, error) {
    expirationTime := time.Now().Add(time.Hour * 1).Unix() // Set expiration time to 1 hour from now

    claims := jwt.MapClaims{
        "adminUUID": adminUUID,
        "email":     email,
        "exp":       expirationTime,
    }

    parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := parseToken.SignedString(secretKey)
    if err != nil {
        return "", err
    }

    return signedToken, nil
}


func VerifyToken(c *gin.Context) (interface{}, error) {
    headerToken := c.Request.Header.Get("Authorization")
    fmt.Println("Header Token:", headerToken)

    bearer := strings.HasPrefix(headerToken, "Bearer")
    fmt.Println("Bearer Prefix Exists:", bearer)

    if !bearer {
        return nil, errors.New(errorSignInToProceed)
    }

    stringToken := strings.Split(headerToken, " ")[1]
    fmt.Println("String Token:", stringToken)

    token, err := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New(errorSignInToProceed)
        }
        return secretKey, nil
    })

    if err != nil || !token.Valid {
        fmt.Println("Token Validation Error:", err)
        return nil, errors.New(errorSignInToProceed)
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        fmt.Println("Failed to extract claims from token.")
        return nil, errors.New(errorSignInToProceed)
    }
    fmt.Println("Claims:", claims)

    expClaim, exists := claims["exp"]
    if !exists {
        return nil, errors.New(errorExpireClaimMissing)
    }

    expTime, ok := expClaim.(float64)
    if !ok {
        return nil, errors.New(errorExpireClaimType)
    }

    if time.Now().Unix() > int64(expTime) {
        return nil, errors.New(errorTokenExpired)
    }

    return token.Claims.(jwt.MapClaims), nil
}
