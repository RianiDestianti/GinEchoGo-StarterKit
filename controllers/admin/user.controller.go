package admin

import (
	"math"
	"net/http"
	"net/url"
	"time"

	"gin-skeleton/config"
	"gin-skeleton/models"
	"gin-skeleton/repositories"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (ctl *UserController) GetData(request *gin.Context) {
	page := repositories.ToInt(request.DefaultQuery("page", "1"), 1)
	limit := repositories.ToInt(request.DefaultQuery("limit", "10"), 10)
	keyword := request.Query("keyword")
	role := request.Query("role")
	offset := (page - 1) * limit

	query := config.DB.Model(&models.UserModel{})

	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", like, like, like)
	}

	if role != "" {
		query = query.Where("role_id = ?", role)
	}

	var total int64
	query.Count(&total)

	var result []models.UserModel
	query.Preload("Role").Order("created_at DESC").Offset(offset).Limit(limit).Find(&result)

	users := []gin.H{}
	for _, item := range result {
		users = append(users, gin.H{
			"id":     item.ID,
			"name":   item.Name,
			"phone":  item.Phone,
			"email":  item.Email,
			"avatar": item.Avatar,
			"role":   item.RoleName(),
		})
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	request.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Success return data user",
		"data": gin.H{
			"data":  users,
			"page":  page,
			"limit": limit,
			"total": gin.H{
				"pages": totalPages,
				"data":  total,
			},
		},
	})
}

func (ctl *UserController) ActionData(request *gin.Context) {
	body := repositories.GetBody(request)
	id := body.Int("id")

	rules := map[string]string{
		"name":     "required",
		"email":    "required|email",
		"role":     "required",
		"password": "required|min:6",
	}
	if id != 0 {
		rules["password"] = "min:6"
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

	var checkEmail models.UserModel
	config.DB.Where("email = ? AND id != ?", body.String("email"), id).Limit(1).Find(&checkEmail)
	if checkEmail.ID != 0 {
		request.JSON(http.StatusOK, gin.H{
			"status":  422,
			"message": "The email has been registered",
			"data": gin.H{
				"email": []string{"The email has been registered"},
			},
		})
		return
	}

	var data models.UserModel
	config.DB.Where("id = ?", id).Limit(1).Find(&data)

	if data.ID != 0 {

		data.Name = body.String("name")
		data.Email = body.String("email")
		data.Phone = body.Nullable("phone")
		data.RoleID = uint(body.Int("role"))
		data.Role = nil

		if body.Has("password") {
			data.Password = repositories.HashPassword(body.String("password"))
		}
	} else {

		now := time.Now()
		avatar := "https://ui-avatars.com/api/?format=png&name=" + url.QueryEscape(body.String("name"))

		data = models.UserModel{
			Name:       body.String("name"),
			Email:      body.String("email"),
			Phone:      body.Nullable("phone"),
			RoleID:     uint(body.Int("role")),
			Password:   repositories.HashPassword(body.String("password")),
			Avatar:     &avatar,
			VerifiedAt: &now,
		}
	}

	if err := config.DB.Save(&data).Error; err != nil {
		request.JSON(http.StatusOK, gin.H{
			"status":  500,
			"message": err.Error(),
		})
		return
	}

	request.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Success saved data",
	})
}

func (ctl *UserController) GetDetail(request *gin.Context) {
	id := request.Param("id")

	var tmpData models.UserModel
	config.DB.Preload("Role").Where("id = ?", id).Limit(1).Find(&tmpData)

	if tmpData.ID != 0 {

		request.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "Success return detail data",
			"data": gin.H{
				"id":          tmpData.ID,
				"name":        tmpData.Name,
				"email":       tmpData.Email,
				"phone":       tmpData.Phone,
				"avatar":      tmpData.Avatar,
				"verified_at": tmpData.VerifiedAt,
				"created_at":  tmpData.CreatedAt,
				"updated_at":  tmpData.UpdatedAt,
				"role_id":     tmpData.RoleID,
				"role":        tmpData.RoleID,
				"roleName":    tmpData.RoleName(),
			},
		})
	} else {
		request.JSON(http.StatusOK, gin.H{
			"status":  404,
			"message": "Data not found",
		})
	}
}

func (ctl *UserController) DeleteData(request *gin.Context) {
	id := request.Param("id")

	var data models.UserModel
	config.DB.Where("id = ?", id).Limit(1).Find(&data)

	if data.ID != 0 {

		config.DB.Delete(&data)

		request.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "Success delete data",
		})
	} else {
		request.JSON(http.StatusOK, gin.H{
			"status":  404,
			"message": "Data not found",
		})
	}
}
