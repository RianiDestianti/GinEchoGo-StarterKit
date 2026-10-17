package models

type VerificationModel struct {
	BaseModel
	Email string `json:"email"`
	Value string `json:"value"`
}

func (VerificationModel) TableName() string {
	return "verifications"
}
