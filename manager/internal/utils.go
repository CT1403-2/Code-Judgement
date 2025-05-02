package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/CT1403-2/Code-Judgement/proto"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
	"os"
	"strings"
	"time"
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var jwtSecret = []byte(GetEnv("JWT_SECRET_KEY", ""))

func GenerateJWT(userId int32, roleType string) (string, error) {
	claims := jwt.MapClaims{
		"userId":   userId,
		"roleType": roleType,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (int32, int32, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return 0, 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		userId := int32(claims["userId"].(float64))
		roleType := proto.Role_value[claims["roleType"].(string)]
		exp := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().After(exp) {
			return 0, 0, errors.New("token is expired")
		}
		return userId, roleType, nil
	}

	return 0, 0, fmt.Errorf("invalid token")
}

func ExtractTokenFromContext(ctx context.Context) (string, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", fmt.Errorf("missing metadata in context")
	}

	// Look for the "authorization" header
	authHeader := md["authorization"]
	if len(authHeader) == 0 {
		return "", "", fmt.Errorf("authorization header not found")
	}

	// Extract the token (assuming "Bearer <token>")
	parts := strings.Split(authHeader[0], " ")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid authorization header format")
	}

	return parts[0], parts[1], nil
}

func IsJudgeServer(token string) bool {
	return token == GetEnv("JUDGE_API_SECRET_KEY", "")
}
