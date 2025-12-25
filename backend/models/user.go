package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	PasswordHash string `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type FailedAttempt struct {
	ID          int       `json:"id"`
	UserID      *int      `json:"user_id"`
	Email       string    `json:"email"`
	AttemptedAt time.Time `json:"attempted_at"`
	IPAddress   string    `json:"ip_address"`
}

