package models

import "time"

type UserModel struct {
	BaseModel
	SoftDelete
	RoleID     uint       `json:"role_id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Phone      *string    `json:"phone"`
	Password   string     `json:"-"`
	Avatar     *string    `json:"avatar"`
	VerifiedAt *time.Time `json:"verified_at"`

	Role *RoleModel `gorm:"foreignKey:RoleID" json:"-"`
}

func (UserModel) TableName() string {
	return "users"
}

func (user UserModel) RoleName() string {
	if user.Role == nil {
		return ""
	}
	return user.Role.Name
}
