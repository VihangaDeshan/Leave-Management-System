package models

import "time"

// LeaveType represents a type of leave (Annual, Sick, etc.)
type LeaveType struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}
