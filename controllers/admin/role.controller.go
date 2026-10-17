package admin

import (
	"net/http"

	"gin-skeleton/config"
	"gin-skeleton/models"

	"github.com/gin-gonic/gin"
)

type RoleController struct{}

func (ctl *RoleController) GetData(request *gin.Context) {
	var roles []models.RoleModel
	config.DB.Find(&roles)

	data := []gin.H{}
	for _, item := range roles {
		data = append(data, gin.H{
			"id":   item.ID,
			"name": item.Name,
		})
	}

	request.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Success return data",
		"data":    data,
	})
}
