package service

import (
	"database/sql"
	"leave-management-system/internal/dto"
	"leave-management-system/internal/repository"
	"leave-management-system/internal/utils"
	apperrors "leave-management-system/pkg/errors"
	"time"
)

type DashboardService struct {
	leaveRepo     *repository.LeaveRepository
	userRepo      *repository.UserRepository
	leaveTypeRepo *repository.LeaveTypeRepository
	db            *sql.DB
}

func NewDashboardService(
	leaveRepo *repository.LeaveRepository,
	userRepo *repository.UserRepository,
	leaveTypeRepo *repository.LeaveTypeRepository,
	db *sql.DB,
) *DashboardService {
	return &DashboardService{
		leaveRepo:     leaveRepo,
		userRepo:      userRepo,
		leaveTypeRepo: leaveTypeRepo,
		db:            db,
	}
}

// GetDashboardStats retrieves dashboard statistics
func (s *DashboardService) GetDashboardStats() (*dto.DashboardStatsResponse, error) {
	stats := &dto.DashboardStatsResponse{}

	// Get pending requests count
	pendingCount, err := s.leaveRepo.GetPendingCount()
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to get pending requests count", err)
	}
	stats.PendingRequests = pendingCount

	// Get total employees count
	totalEmployees, err := s.userRepo.GetActiveEmployeeCount()
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to get employee count", err)
	}
	stats.TotalEmployees = totalEmployees

	// Get employees on leave today
	today := time.Now()
	leavesToday, err := s.leaveRepo.GetEmployeesOnLeave(today)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to get employees on leave", err)
	}
	stats.EmployeesOnLeave = len(leavesToday)

	// Get detailed info about employees on leave
	var employeesOnLeave []dto.UserOnLeave
	for _, leave := range leavesToday {
		user, err := s.userRepo.FindByID(leave.UserID)
		if err != nil {
			continue
		}

		leaveType, err := s.leaveTypeRepo.FindByID(leave.LeaveTypeID)
		if err != nil {
			continue
		}

		employeesOnLeave = append(employeesOnLeave, dto.UserOnLeave{
			UserID:         user.ID,
			FullName:       user.GetFullName(),
			Email:          user.Email,
			Department:     user.Department,
			LeaveType:      leaveType.Name,
			StartDate:      utils.FormatDateString(leave.StartDate),
			EndDate:        utils.FormatDateString(leave.EndDate),
			TotalDays:      leave.TotalDays,
			LeaveRequestID: leave.ID,
		})
	}
	stats.EmployeesOnLeaveToday = employeesOnLeave

	// Get leave statistics by type
	leavesByType, err := s.getLeaveStatsByType()
	if err != nil {
		return nil, err
	}
	stats.LeavesByType = leavesByType

	// Calculate approved and rejected counts from leavesByType
	approvedTotal := 0
	rejectedTotal := 0
	for _, stat := range leavesByType {
		approvedTotal += stat.ApprovedCount
		rejectedTotal += stat.RejectedCount
	}
	stats.ApprovedRequests = approvedTotal
	stats.RejectedRequests = rejectedTotal

	// Get recent leave requests (last 5)
	recentLeaves, err := s.getRecentLeaveRequests(5)
	if err != nil {
		return nil, err
	}
	stats.RecentLeaveRequests = recentLeaves

	return stats, nil
}

// getLeaveStatsByType calculates statistics grouped by leave type
func (s *DashboardService) getLeaveStatsByType() ([]dto.LeaveTypeStats, error) {
	query := `
		SELECT
			lt.name as leave_type_name,
			COUNT(*) as total_requests,
			SUM(CASE WHEN lr.status = 'approved' THEN 1 ELSE 0 END) as approved_count,
			SUM(CASE WHEN lr.status = 'rejected' THEN 1 ELSE 0 END) as rejected_count,
			SUM(CASE WHEN lr.status = 'pending' THEN 1 ELSE 0 END) as pending_count,
			SUM(CASE WHEN lr.status = 'approved' THEN lr.total_days ELSE 0 END) as total_days
		FROM leave_requests lr
		JOIN leave_types lt ON lr.leave_type_id = lt.id
		GROUP BY lt.id, lt.name
		ORDER BY total_requests DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to get leave statistics", err)
	}
	defer rows.Close()

	var stats []dto.LeaveTypeStats
	for rows.Next() {
		var stat dto.LeaveTypeStats
		err := rows.Scan(
			&stat.LeaveTypeName,
			&stat.TotalRequests,
			&stat.ApprovedCount,
			&stat.RejectedCount,
			&stat.PendingCount,
			&stat.TotalDays,
		)
		if err != nil {
			continue
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// getRecentLeaveRequests retrieves the most recent leave requests
func (s *DashboardService) getRecentLeaveRequests(limit int) ([]dto.LeaveRequestResponse, error) {
	query := `
		SELECT lr.id, lr.user_id, lr.leave_type_id, lr.start_date, lr.end_date,
		       lr.total_days, lr.reason, lr.status, lr.reviewed_by, lr.reviewed_at,
		       lr.review_notes, lr.created_at, lr.updated_at,
		       u.first_name, u.last_name, u.email,
		       lt.name as leave_type_name
		FROM leave_requests lr
		JOIN users u ON lr.user_id = u.id
		JOIN leave_types lt ON lr.leave_type_id = lt.id
		ORDER BY lr.created_at DESC
		LIMIT $1
	`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to get recent leave requests", err)
	}
	defer rows.Close()

	var leaves []dto.LeaveRequestResponse
	for rows.Next() {
		var leave dto.LeaveRequestResponse
		var user dto.UserBasicInfo
		var leaveTypeName string
		var startDate, endDate, createdAt, updatedAt time.Time
		var reviewedAt sql.NullTime

		err := rows.Scan(
			&leave.ID,
			&leave.UserID,
			&leave.LeaveTypeID,
			&startDate,
			&endDate,
			&leave.TotalDays,
			&leave.Reason,
			&leave.Status,
			&leave.ReviewedBy,
			&reviewedAt,
			&leave.ReviewNotes,
			&createdAt,
			&updatedAt,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&leaveTypeName,
		)
		if err != nil {
			continue
		}

		// Format dates
		leave.StartDate = utils.FormatDateString(startDate)
		leave.EndDate = utils.FormatDateString(endDate)
		leave.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		leave.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")

		// Set user info
		user.ID = leave.UserID
		user.FullName = user.FirstName + " " + user.LastName
		leave.User = &user

		// Set leave type
		leave.LeaveType = &dto.LeaveTypeResponse{
			ID:   leave.LeaveTypeID,
			Name: leaveTypeName,
		}

		// Handle reviewed_at
		if reviewedAt.Valid {
			reviewedAtStr := reviewedAt.Time.Format("2006-01-02 15:04:05")
			leave.ReviewedAt = &reviewedAtStr
		}

		leaves = append(leaves, leave)
	}

	return leaves, nil
}
