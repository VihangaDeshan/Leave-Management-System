package models

import "time"

// LeaveBalance represents a user's leave balance for a specific leave type and year
type LeaveBalance struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	LeaveTypeID   int       `json:"leave_type_id"`
	TotalDays     float64   `json:"total_days"`
	UsedDays      float64   `json:"used_days"`
	AvailableDays float64   `json:"available_days"` // Computed field
	Year          int       `json:"year"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relationships (populated by joins)
	LeaveType *LeaveType `json:"leave_type,omitempty"`
}

// HasSufficientBalance checks if enough leave days are available
func (lb *LeaveBalance) HasSufficientBalance(requestedDays float64) bool {
	return lb.AvailableDays >= requestedDays
}
