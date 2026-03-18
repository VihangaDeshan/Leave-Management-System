package dto

// CreateLeaveRequest represents a request to create a leave request
type CreateLeaveRequest struct {
	LeaveTypeID int    `json:"leave_type_id" validate:"required"`
	StartDate   string `json:"start_date" validate:"required"` // Format: YYYY-MM-DD
	EndDate     string `json:"end_date" validate:"required"`   // Format: YYYY-MM-DD
	Reason      string `json:"reason" validate:"required,min=10"`
}

// UpdateLeaveRequest represents a request to update a leave request
type UpdateLeaveRequest struct {
	LeaveTypeID int    `json:"leave_type_id"`
	StartDate   string `json:"start_date"` // Format: YYYY-MM-DD
	EndDate     string `json:"end_date"`   // Format: YYYY-MM-DD
	Reason      string `json:"reason"`
}

// ReviewLeaveRequest represents a request to approve/reject a leave
type ReviewLeaveRequest struct {
	Status      string  `json:"status" validate:"required,oneof=approved rejected"`
	ReviewNotes *string `json:"review_notes"`
}

// LeaveRequestResponse represents a leave request in API responses
type LeaveRequestResponse struct {
	ID          int                   `json:"id"`
	UserID      int                   `json:"user_id"`
	User        *UserBasicInfo        `json:"user,omitempty"`
	LeaveTypeID int                   `json:"leave_type_id"`
	LeaveType   *LeaveTypeResponse    `json:"leave_type,omitempty"`
	StartDate   string                `json:"start_date"`
	EndDate     string                `json:"end_date"`
	TotalDays   float64               `json:"total_days"`
	Reason      string                `json:"reason"`
	Status      string                `json:"status"`
	ReviewedBy  *int                  `json:"reviewed_by"`
	Reviewer    *UserBasicInfo        `json:"reviewer,omitempty"`
	ReviewedAt  *string               `json:"reviewed_at"`
	ReviewNotes *string               `json:"review_notes"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
}

// LeaveRequestsListResponse represents a paginated list of leave requests
type LeaveRequestsListResponse struct {
	LeaveRequests []LeaveRequestResponse `json:"leave_requests"`
	TotalCount    int                    `json:"total_count"`
	Page          int                    `json:"page"`
	PageSize      int                    `json:"page_size"`
	TotalPages    int                    `json:"total_pages"`
}

// UserBasicInfo represents basic user information for nested responses
type UserBasicInfo struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
}

// LeaveTypeResponse represents a leave type in API responses
type LeaveTypeResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IsActive    bool    `json:"is_active"`
}

// LeaveBalanceResponse represents a leave balance in API responses
type LeaveBalanceResponse struct {
	ID            int                `json:"id"`
	UserID        int                `json:"user_id"`
	LeaveTypeID   int                `json:"leave_type_id"`
	LeaveType     *LeaveTypeResponse `json:"leave_type,omitempty"`
	TotalDays     float64            `json:"total_days"`
	UsedDays      float64            `json:"used_days"`
	AvailableDays float64            `json:"available_days"`
	Year          int                `json:"year"`
}

// AllocateBalanceRequest represents a request to allocate leave balance
type AllocateBalanceRequest struct {
	UserID      int     `json:"user_id" validate:"required"`
	LeaveTypeID int     `json:"leave_type_id" validate:"required"`
	TotalDays   float64 `json:"total_days" validate:"required,gt=0"`
	Year        int     `json:"year" validate:"required"`
}

// UpdateBalanceRequest represents a request to update leave balance
type UpdateBalanceRequest struct {
	TotalDays float64 `json:"total_days" validate:"required,gt=0"`
}
