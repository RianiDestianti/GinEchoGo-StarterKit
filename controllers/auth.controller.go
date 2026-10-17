package controllers

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"gin-skeleton/config"
	"gin-skeleton/middleware"
	"gin-skeleton/models"
	"gin-skeleton/notifications"
	"gin-skeleton/repositories"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct{}

func (ctl *AuthController) Login(request *gin.Context) {
	body := repositories.GetBody(request)

	rules := map[string]string{
		"email":    "required|email",
		"password": "required",
	}

	validation := repositories.NewValidator(body, rules)
	if validation.Fails() {
		request.JSON(http.StatusOK, gin.H{
			"status":  422,
			"message": "Validation failed",
			"data":    validation.Errors.All(),
		})
		return
	}

	var user models.UserModel
	config.DB.Preload("Role").Where("email = ?", body.String("email")).Limit(1).Find(&user)

	if user.ID == 0 {
		request.JSON(http.StatusOK, gin.H{
			"status":  401,
			"message": "Invalid credentials",
		})
		return
	}

	if body.Has("role") && body.String("role") != user.RoleName() {
		request.JSON(http.StatusOK, gin.H{
			"status":  401,
			"message": "Invalid credentials",
		})
		return
	}

	if !repositories.ComparePassword(body.String("password"), user.Password) {
		request.JSON(http.StatusOK, gin.H{
			"status":  401,
			"message": "Invalid credentials",
		})
		return
	}

	token, _ := middleware.Sign(jwt.MapClaims{"email": user.Email, "role": user.RoleName()})
	request.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Success logined",
		"data": gin.H{
			"key":   token,
			"email": user.Email,
		},
	})
}

func (ctl *AuthController) Register(request *gin.Context) {
	payload := repositories.GetBody(request)

	rules := map[string]string{
		"name":     "required",
		"email":    "required|email",
		"password": "required|min:6",
	}

	validation := repositories.NewValidator(payload, rules)
	if validation.Fails() {
		request.JSON(http.StatusOK, gin.H{
			"status":  422,
			"message": "Validation failed",
			"data":    validation.Errors.All(),
		})
		return
	}

	var checkUser models.UserModel
	config.DB.Where("email = ?", payload.String("email")).Limit(1).Find(&checkUser)

	if checkUser.ID != 0 {
		request.JSON(http.StatusOK, gin.H{
			"status":  422,
			"message": "The email has been registered",
			"data": gin.H{
				"email": []string{"The email has been registered"},
			},
		})
		return
	}

	avatar := "https://ui-avatars.com/api/?format=png&name=" + url.QueryEscape(payload.String("name"))
	user := models.UserModel{
		Email:    payload.String("email"),
		Name:     payload.String("name"),
		Password: repositories.HashPassword(payload.String("password")),
		RoleID:   3,
		Avatar:   &avatar,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		request.JSON(http.StatusOK, gin.H{
			"status":  500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	token := repositories.GenerateNumber(6)
	config.DB.Create(&models.VerificationModel{
		Email: payload.String("email"),
		Value: token,
	})

	notifications.RegistrationNotification.Send(payload.String("email"), payload.String("name"), token)
	register, _ := middleware.Sign(jwt.MapClaims{"email": payload.String("email"), "role": "user"})

	request.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Successfully register",
		"data": gin.H{
			"key":   register,
			"email": payload.String("email"),
		},
	})
}

func (ctl *AuthController) Verify(request *gin.Context) {
	body := repositories.GetBody(request)

	rules := map[string]string{
		"email": "required|email",
		"code":  "required",
	}

	validation := repositories.NewValidator(body, rules)
	if validation.Fails() {
		request.JSON(http.StatusOK, gin.H{
			"status":  422,
			"message": "Validation failed",
			"data":    validation.Errors.All(),
		})
		return
	}

	var user models.UserModel
	config.DB.Where("email = ?", body.String("email")).Limit(1).Find(&user)
	if user.ID == 0 {
		request.JSON(http.StatusOK, gin.H{
			"status":  422,
			"message": "Invalid code",
		})
		return
	}

	var getVerification models.VerificationModel
	config.DB.Where("email = ? AND value = ?", body.String("email"), body.String("code")).Limit(1).Find(&getVerification)
	if getVerification.ID == 0 {
		request.JSON(http.StatusOK, gin.H{
			"status":  422,
			"message": "Invalid code",
		})
		return
	}

	now := time.Now()
	user.VerifiedAt = &now
	config.DB.Save(&user)

	config.DB.Where("email = ?", body.String("email")).Delete(&models.VerificationModel{})

	request.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Successfully verification",
	})
}

func (ctl *AuthController) RetryCode(request *gin.Context) {
	payload := repositories.GetBody(request)

	rules := map[string]string{
		"email": "required|email",
	}

	validation := repositories.NewValidator(payload, rules)
	if validation.Fails() {
		request.JSON(http.StatusOK, gin.H{
			"status":  422,
			"message": "Validation failed",
			"data":    validation.Errors.All(),
		})
		return
	}

	var user models.UserModel
	config.DB.Where("email = ?", payload.String("email")).Limit(1).Find(&user)

	config.DB.Where("email = ?", payload.String("email")).Delete(&models.VerificationModel{})

	token := repositories.GenerateNumber(6)
	config.DB.Create(&models.VerificationModel{
		Email: payload.String("email"),
		Value: token,
	})

	name := "User"
	if user.ID != 0 {
		name = user.Name
	}

	notifications.RetryNotification.Send(payload.String("email"), name, token)

	request.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Successfully send code",
	})
}

func (ctl *AuthController) ResetPassword(request *gin.Context) {
	body := repositories.GetBody(request)

	rules := map[string]string{
		"email":    "required|email",
		"password": "required|min:6",
	}

	validation := repositories.NewValidator(body, rules)
	if validation.Fails() {
		request.JSON(http.StatusOK, gin.H{
			"status":  422,
			"message": "Validation failed",
			"data":    validation.Errors.All(),
		})
		return
	}

	email := strings.TrimSpace(body.String("email"))

	var user models.UserModel
	config.DB.Where("email = ?", email).Limit(1).Find(&user)

	if user.ID == 0 {
		request.JSON(http.StatusOK, gin.H{
			"status":  422,
			"message": "Email tidak ditemukan",
			"data": gin.H{
				"email": []string{"Email tidak ditemukan."},
			},
		})
		return
	}

	user.Password = repositories.HashPassword(body.String("password"))
	config.DB.Save(&user)

	request.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Password berhasil direset",
	})
}
