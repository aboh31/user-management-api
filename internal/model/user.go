package model

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Uuid      string `gorm:"Uuid"`
	Username  string `gorm:"unique"`
	Email     string
	Password  string
	CreatedAt time.Time
}

type UserResponse struct {
	Id        string `json:"id"` //ini buat uuid
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func (d *User) ConvertToResponse() UserResponse {
	return UserResponse{
		Id:        d.Uuid,
		Username:  d.Username,
		Email:     d.Email,
		CreatedAt: d.CreatedAt.Format("2006-01-02T15:04:05"),
	}
}
