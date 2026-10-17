package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gin-skeleton/config"
	"gin-skeleton/database"
	"gin-skeleton/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	requiredEnv := []string{"DB_HOST", "DB_USER", "DB_NAME", "JWT_SECRET"}
	if os.Getenv("DB_CLIENT") == "sqlite" {
		requiredEnv = []string{"DB_NAME", "JWT_SECRET"}
	}

	missing := []string{}
	for _, key := range requiredEnv {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		fmt.Println("[error] Env wajib belum diisi: " + strings.Join(missing, ", "))
		os.Exit(1)
	}

	config.ConnectDatabase()

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		database.Seed()
		return
	}

	if os.Getenv("APP_DEBUG") != "true" {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()
	app.Use(gin.Logger())
	app.Use(gin.CustomRecovery(func(request *gin.Context, err any) {
		log.Println("[error] Unhandled error:", err)
		request.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": fmt.Sprint(err),
			"data":    nil,
		})
	}))

	corsConfig := cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}
	if origin := os.Getenv("CORS_ORIGIN"); origin == "" || origin == "*" {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = strings.Split(origin, ",")
	}
	app.Use(cors.New(corsConfig))

	routes.AdminRoutes(app)
	routes.UserRoutes(app)

	app.HandleMethodNotAllowed = true
	app.NoRoute(func(request *gin.Context) {
		request.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": fmt.Sprintf("Cannot %s %s", request.Request.Method, request.Request.URL.Path),
			"data":    nil,
		})
	})
	app.NoMethod(func(request *gin.Context) {
		request.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":  405,
			"message": "The method is not allowed for the requested URL.",
			"data":    nil,
		})
	})

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on http://%s:%s", host, port)
	if err := app.Run(host + ":" + port); err != nil {
		log.Fatal(err)
	}
}
