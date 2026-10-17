package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SoftDelete: db.Delete() hanya mengisi deleted_at, dan semua query otomatis
// mengabaikan data yang sudah dihapus. Pakai db.Unscoped() untuk ikut mengambilnya.
type SoftDelete struct {
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
