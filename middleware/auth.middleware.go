package middleware

import (
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gin-skeleton/config"
	"gin-skeleton/models"
	"gin-skeleton/types"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var units = map[string]time.Duration{
	"s": time.Second, "m": time.Minute, "h": time.Hour,
	"d": 24 * time.Hour, "w": 7 * 24 * time.Hour, "y": 8766 * time.Hour,
}

func Sign(claims jwt.MapClaims) (string, error) {
	expiresIn := time.Duration(8766) * time.Hour

	if match := regexp.MustCompile(`^(\d+)\s*([smhdwy]?)$`).FindStringSubmatch(strings.TrimSpace(os.Getenv("JWT_EXPIRES_IN"))); match != nil {
		value, _ := strconv.Atoi(match[1])
		unit := match[2]
		if unit == "" {
			unit = "s"
		}
		expiresIn = time.Duration(value) * units[unit]
	}

	now := time.Now()
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(expiresIn).Unix()

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func AuthMiddleware(request *gin.Context) {
	authHeader := request.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		request.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Unauthorized", "data": nil})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	secret := os.Getenv("JWT_SECRET")

	if token == "" || secret == "" {
		request.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Missing token or JWT secret", "data": nil})
		return
	}

	claims := jwt.MapClaims{}
	decoded, err := jwt.ParseWithClaims(token, claims, func(parsedToken *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	var user models.UserModel
	if err == nil && decoded.Valid {
		email, _ := claims["email"].(string)
		config.DB.Preload("Role").Where("email = ?", email).Limit(1).Find(&user)
	}

	if user.ID == 0 {
		request.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": 403, "message": "Invalid or expired token", "data": nil})
		return
	}

	request.Set("user", types.AuthUser{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.RoleName(),
		Avatar: user.Avatar,
	})
	request.Next()
}

func AdminOnly(request *gin.Context) {
	user, ok := request.Get("user")

	if !ok || user.(types.AuthUser).Role != "admin" {
		request.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": 403, "message": "Forbidden: Admins only", "data": nil})
		return
	}

	request.Next()
}

func TutornAdmin(request *gin.Context) {
	user, ok := request.Get("user")

	if !ok || user.(types.AuthUser).Role != "admin" && user.(types.AuthUser).Role != "tutor" {
		request.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": 403, "message": "Forbidden: Admins and Tutor only", "data": nil})
		return
	}

	request.Next()
}
