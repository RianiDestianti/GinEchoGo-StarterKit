package database

import (
	"log"
	"os"
	"time"

	"gin-skeleton/config"
	"gin-skeleton/models"
	"gin-skeleton/repositories"
)

// Seed membuat akun admin dari ADMIN_* di .env
//
//	go run . seed
func Seed() {
	adminName := os.Getenv("ADMIN_NAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminName == "" {
		adminName = "Administrator"
	}

	if adminEmail == "" || adminPassword == "" {
		log.Fatal("[error] ADMIN_EMAIL dan ADMIN_PASSWORD wajib diisi di .env")
	}

	var exists int64
	config.DB.Unscoped().Model(&models.UserModel{}).Where("email = ?", adminEmail).Count(&exists)

	if exists > 0 {
		log.Printf("[info] Admin %s sudah ada, seed dilewati", adminEmail)
		return
	}

	now := time.Now()
	err := config.DB.Create(&models.UserModel{
		RoleID:     1,
		Name:       adminName,
		Email:      adminEmail,
		Password:   repositories.HashPassword(adminPassword),
		VerifiedAt: &now,
	}).Error

	if err != nil {
		log.Fatalf("[error] Gagal seed: %v", err)
	}

	log.Printf("[ok] Admin %s berhasil dibuat", adminEmail)
}
