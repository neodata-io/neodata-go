package dto

// UserLoginDTO represents the data required to log in a user.
type UserLoginDTO struct {
	Email    string `json:"email" validate:"required,email"`    // Email is required and must be a valid email format
	Password string `json:"password" validate:"required,min=8"` // Password is required with a minimum length of 8
}
