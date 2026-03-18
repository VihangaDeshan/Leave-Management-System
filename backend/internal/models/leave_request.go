package models

import "time"

// LeaveRequest represents a leave request submitted by a user
type LeaveRequest struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	LeaveTypeID int        `json:"leave_type_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	TotalDays   float64    `json:"total_days"`
	Reason      string     `json:"reason"`
	Status      string     `json:"status"` // pending, approved, rejected, cancelled
	ReviewedBy  *int       `json:"reviewed_by"`
	ReviewedAt  *time.Time `json:"reviewed_at"`
	ReviewNotes *string    `json:"review_notes"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relationships (populated by joins)
	User      *User      `json:"user,omitempty"`
	LeaveType *LeaveType `json:"leave_type,omitempty"`
	Reviewer  *User      `json:"reviewer,omitempty"`
}

// IsPending checks if the leave request is pending
func (lr *LeaveRequest) IsPending() bool {
	return lr.Status == "pending"
}

// IsApproved checks if the leave request is approved
func (lr *LeaveRequest) IsApproved() bool {
	return lr.Status == "approved"
}

// IsRejected checks if the leave request is rejected
func (lr *LeaveRequest) IsRejected() bool {
	return lr.Status == "rejected"
}

// IsCancelled checks if the leave request is cancelled
func (lr *LeaveRequest) IsCancelled() bool {
	return lr.Status == "cancelled"
}

// CanBeModified checks if the leave request can be modified or cancelled
func (lr *LeaveRequest) CanBeModified() bool {
	return lr.Status == "pending"
}
