package response

import "time"

// UserResponse represents user data in response
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthResponse represents authentication response with token
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
