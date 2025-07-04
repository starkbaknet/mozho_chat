package dto

type UpdateUserRequest struct {
    Username *string          `json:"username,omitempty"`
    Email    *string          `json:"email,omitempty"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=6"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}
