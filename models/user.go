package models

import (
	"strings"
	"time"
)

type (
	User struct {
		ID        int       `gorm:"column:id"`
		Username  string    `gorm:"column:username;size:50;not null"`
		Email     string    `gorm:"column:email;size:50;not null;unique"`
		Password  string    `gorm:"column:password;size:255;not null"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
	}

	UserRegisterRequest struct {
		Username string `json:"username" valid:"type(string),required,length(1|50)"`
		Email    string `json:"email" valid:"email,required,length(3|50)"`
		Password string `json:"password" valid:"type(string),required,length(6|255)"`
	}

	UserLoginRequest struct {
		Email    string `json:"email" valid:"email,required,length(3|50)"`
		Password string `json:"password" valid:"type(string),required,length(6|255)"`
	}

	UserUpdateRequest struct {
		Username string `json:"username"`
	}

	UserRegisterResponse struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	UserResponse struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	UserLoginResponse struct {
		ID          int    `json:"id"`
		Token       string `json:"token"`
		TokenExpire string `json:"expire"`
	}
)

// override gorm table name
func (User) TableName() string {
	return "users"
}

func (req *UserRegisterRequest) Trim() {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
}

func (req *UserLoginRequest) Trim() {
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
}

func (req *UserUpdateRequest) Trim() {
	req.Username = strings.TrimSpace(req.Username)
}
