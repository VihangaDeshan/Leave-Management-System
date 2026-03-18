package dto

// DashboardStatsResponse represents dashboard statistics
type DashboardStatsResponse struct {
	PendingRequests       int                    `json:"pending_requests"`
	ApprovedRequests      int                    `json:"approved_requests"`
	RejectedRequests      int                    `json:"rejected_requests"`
	EmployeesOnLeave      int                    `json:"employees_on_leave"`
	TotalEmployees        int                    `json:"total_employees"`
	LeavesByType          []LeaveTypeStats       `json:"leaves_by_type"`
	RecentLeaveRequests   []LeaveRequestResponse `json:"recent_leave_requests"`
	EmployeesOnLeaveToday []UserOnLeave          `json:"employees_on_leave_today"`
}

// LeaveTypeStats represents leave statistics by type
type LeaveTypeStats struct {
	LeaveTypeName string  `json:"leave_type_name"`
	TotalRequests int     `json:"total_requests"`
	ApprovedCount int     `json:"approved_count"`
	RejectedCount int     `json:"rejected_count"`
	PendingCount  int     `json:"pending_count"`
	TotalDays     float64 `json:"total_days"`
}

// UserOnLeave represents a user currently on leave
type UserOnLeave struct {
	UserID         int     `json:"user_id"`
	FullName       string  `json:"full_name"`
	Email          string  `json:"email"`
	Department     *string `json:"department"`
	LeaveType      string  `json:"leave_type"`
	StartDate      string  `json:"start_date"`
	EndDate        string  `json:"end_date"`
	TotalDays      float64 `json:"total_days"`
	LeaveRequestID int     `json:"leave_request_id"`
}
