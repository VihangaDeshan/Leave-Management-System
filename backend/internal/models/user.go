package models

import "time"

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never send password hash in JSON
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Role         string    `json:"role"` // employee, manager, admin
	Department   *string   `json:"department"`
	ManagerID    *int      `json:"manager_id"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsEmployee checks if the user has employee role
func (u *User) IsEmployee() bool {
	return u.Role == "employee"
}

// IsManager checks if the user has manager role
func (u *User) IsManager() bool {
	return u.Role == "manager"
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// HasAdminAccess checks if the user has admin or manager privileges
func (u *User) HasAdminAccess() bool {
	return u.Role == "admin" || u.Role == "manager"
}
