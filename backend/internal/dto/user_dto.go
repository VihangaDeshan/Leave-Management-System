package dto

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Email      string  `json:"email" validate:"required,email"`
	Password   string  `json:"password" validate:"required,min=8"`
	FirstName  string  `json:"first_name" validate:"required"`
	LastName   string  `json:"last_name" validate:"required"`
	Role       string  `json:"role" validate:"required,oneof=employee manager admin"`
	Department *string `json:"department"`
	ManagerID  *int    `json:"manager_id"`
}

// UpdateUserRequest represents a request to update user information
type UpdateUserRequest struct {
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	Department *string `json:"department"`
	ManagerID  *int    `json:"manager_id"`
}

// UpdateProfileRequest represents a request to update own profile
type UpdateProfileRequest struct {
	FirstName  string  `json:"first_name" validate:"required"`
	LastName   string  `json:"last_name" validate:"required"`
	Department *string `json:"department"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID         int     `json:"id"`
	Email      string  `json:"email"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	FullName   string  `json:"full_name"`
	Role       string  `json:"role"`
	Department *string `json:"department"`
	ManagerID  *int    `json:"manager_id"`
	IsActive   bool    `json:"is_active"`
	CreatedAt  string  `json:"created_at"`
}

// UsersListResponse represents a paginated list of users
type UsersListResponse struct {
	Users      []UserResponse `json:"users"`
	TotalCount int            `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}
