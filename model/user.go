package model

type User struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}
