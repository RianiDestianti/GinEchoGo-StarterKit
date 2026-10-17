package user

import (
	"net/http"

	"gin-skeleton/config"
	"gin-skeleton/models"
	"gin-skeleton/repositories"
	"gin-skeleton/types"

	"github.com/gin-gonic/gin"
)

type ProfileController struct{}

func (ctl *ProfileController) GetProfile(request *gin.Context) {
	user := request.MustGet("user").(types.AuthUser)

	var getUser models.UserModel
	config.DB.Where("id = ?", user.ID).Limit(1).Find(&getUser)

	request.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Success return data user",
		"data": gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"role":     user.Role,
			"avatar":   user.Avatar,
			"verified": getUser.ID != 0 && getUser.VerifiedAt != nil,
		},
	})
}

func (ctl *ProfileController) Action(request *gin.Context) {
	body := repositories.GetBody(request)
	authUser := request.MustGet("user").(types.AuthUser)

	rules := map[string]string{
		"name":     "required",
		"password": "min:6",
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
	config.DB.Where("id = ?", authUser.ID).Limit(1).Find(&user)
	if user.ID == 0 {
		request.JSON(http.StatusOK, gin.H{
			"status":  404,
			"message": "Data not found",
		})
		return
	}

	user.Name = body.String("name")

	if _, ok := body["phone"]; ok {
		user.Phone = body.Nullable("phone")
	}

	if _, ok := body["avatar"]; ok {
		user.Avatar = body.Nullable("avatar")
	}

	if body.Has("password") {
		user.Password = repositories.HashPassword(body.String("password"))
	}

	if err := config.DB.Save(&user).Error; err != nil {
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
