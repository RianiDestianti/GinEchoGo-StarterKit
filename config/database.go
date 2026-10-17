package config

import (
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	var dialector gorm.Dialector

	if os.Getenv("DB_CLIENT") == "sqlite" {
		dialector = sqlite.Open(os.Getenv("DB_NAME"))
	} else {
		cfg := mysqldriver.NewConfig()
		cfg.User = os.Getenv("DB_USER")
		cfg.Passwd = os.Getenv("DB_PASSWORD")
		cfg.Net = "tcp"
		cfg.Addr = os.Getenv("DB_HOST") + ":" + getenv("DB_PORT", "3306")
		cfg.DBName = os.Getenv("DB_NAME")
		cfg.ParseTime = true
		cfg.Loc = time.Local
		cfg.Params = map[string]string{"charset": "utf8mb4"}

		dialector = mysql.Open(cfg.FormatDSN())
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
		}),
	})
	if err != nil {
		log.Fatalf("[error] Gagal koneksi database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[error] Gagal koneksi database: %v", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxIdleTime(30 * time.Second)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	DB = db
}

func getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
