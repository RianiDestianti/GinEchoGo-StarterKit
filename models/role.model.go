package models

type RoleModel struct {
	BaseModel
	SoftDelete
	Name string `json:"name"`
}

func (RoleModel) TableName() string {
	return "roles"
}
