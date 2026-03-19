package service

import (
	"leave-management-system/internal/dto"
	"leave-management-system/internal/models"
	"leave-management-system/internal/repository"
	"leave-management-system/internal/utils"
	apperrors "leave-management-system/pkg/errors"
	"math"
	"time"
)

type LeaveService struct {
	leaveRepo     *repository.LeaveRepository
	balanceRepo   *repository.BalanceRepository
	leaveTypeRepo *repository.LeaveTypeRepository
	userRepo      *repository.UserRepository
}

func NewLeaveService(
	leaveRepo *repository.LeaveRepository,
	balanceRepo *repository.BalanceRepository,
	leaveTypeRepo *repository.LeaveTypeRepository,
	userRepo *repository.UserRepository,
) *LeaveService {
	return &LeaveService{
		leaveRepo:     leaveRepo,
		balanceRepo:   balanceRepo,
		leaveTypeRepo: leaveTypeRepo,
		userRepo:      userRepo,
	}
}

// CreateLeaveRequest creates a new leave request
func (s *LeaveService) CreateLeaveRequest(userID int, req *dto.CreateLeaveRequest) (*dto.LeaveRequestResponse, error) {
	// Parse dates
	startDate, err := utils.ParseDateString(req.StartDate)
	if err != nil {
		return nil, apperrors.BadRequest("Invalid start date format. Use YYYY-MM-DD")
	}

	endDate, err := utils.ParseDateString(req.EndDate)
	if err != nil {
		return nil, apperrors.BadRequest("Invalid end date format. Use YYYY-MM-DD")
	}

	// Validate date range
	if !utils.IsDateRangeValid(startDate, endDate) {
		return nil, apperrors.BadRequest("End date must be after or equal to start date")
	}

	// Check if dates are in the past
	if startDate.Before(utils.GetStartOfDay(time.Now())) {
		return nil, apperrors.BadRequest("Cannot request leave for past dates")
	}

	// Calculate working days
	totalDays := utils.CalculateWorkingDays(startDate, endDate)
	if totalDays == 0 {
		return nil, apperrors.BadRequest("Leave request must include at least one working day")
	}

	// Check for overlapping leave requests
	hasOverlap, err := s.leaveRepo.CheckOverlap(userID, startDate, endDate, nil)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to check leave overlap", err)
	}
	if hasOverlap {
		return nil, apperrors.Conflict("You already have a leave request for this period")
	}

	// Check leave balance
	currentYear := utils.GetCurrentYear()
	if err := s.ensureUserHasCurrentYearBalances(userID, currentYear); err != nil {
		return nil, err
	}

	balance, err := s.balanceRepo.FindByUserAndType(userID, req.LeaveTypeID, currentYear)
	if err != nil {
		return nil, apperrors.BadRequest("Leave balance not found for this leave type")
	}

	if !balance.HasSufficientBalance(totalDays) {
		return nil, apperrors.BadRequest("Insufficient leave balance. Available: " +
			utils.FormatDateString(time.Time{}) + " days")
	}

	// Create leave request
	leave := &models.LeaveRequest{
		UserID:      userID,
		LeaveTypeID: req.LeaveTypeID,
		StartDate:   startDate,
		EndDate:     endDate,
		TotalDays:   totalDays,
		Reason:      req.Reason,
		Status:      "pending",
	}

	if err := s.leaveRepo.Create(leave); err != nil {
		return nil, apperrors.InternalServerError("Failed to create leave request", err)
	}

	return s.mapToLeaveRequestResponse(leave), nil
}

// GetUserLeaveRequests retrieves all leave requests for a user
func (s *LeaveService) GetUserLeaveRequests(userID, page, pageSize int) (*dto.LeaveRequestsListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	leaves, totalCount, err := s.leaveRepo.FindByUserID(userID, pageSize, offset)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to retrieve leave requests", err)
	}

	// Enrich with related data
	var leaveResponses []dto.LeaveRequestResponse
	for _, leave := range leaves {
		enrichedLeave, err := s.enrichLeaveRequest(leave)
		if err != nil {
			// Log error but continue
			leaveResponses = append(leaveResponses, *s.mapToLeaveRequestResponse(leave))
			continue
		}
		leaveResponses = append(leaveResponses, *enrichedLeave)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &dto.LeaveRequestsListResponse{
		LeaveRequests: leaveResponses,
		TotalCount:    totalCount,
		Page:          page,
		PageSize:      pageSize,
		TotalPages:    totalPages,
	}, nil
}

// GetLeaveRequestByID retrieves a specific leave request
func (s *LeaveService) GetLeaveRequestByID(id, userID int, isAdmin bool) (*dto.LeaveRequestResponse, error) {
	leave, err := s.leaveRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFound("Leave request not found")
	}

	// Check authorization (user can only view their own, admin can view all)
	if !isAdmin && leave.UserID != userID {
		return nil, apperrors.Forbidden("Not authorized to view this leave request")
	}

	return s.enrichLeaveRequest(leave)
}

// CancelLeaveRequest cancels a pending leave request
func (s *LeaveService) CancelLeaveRequest(id, userID int, isAdmin bool) error {
	leave, err := s.leaveRepo.FindByID(id)
	if err != nil {
		return apperrors.NotFound("Leave request not found")
	}

	// Check authorization
	if !isAdmin && leave.UserID != userID {
		return apperrors.Forbidden("Not authorized to cancel this leave request")
	}

	// Can only cancel pending requests
	if !leave.CanBeModified() {
		return apperrors.BadRequest("Only pending leave requests can be cancelled")
	}

	if err := s.leaveRepo.Delete(id); err != nil {
		return apperrors.InternalServerError("Failed to cancel leave request", err)
	}

	return nil
}

// GetAllLeaveRequests retrieves all leave requests (admin only)
func (s *LeaveService) GetAllLeaveRequests(requesterID int, requesterRole, status string, page, pageSize int) (*dto.LeaveRequestsListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var (
		leaves      []*models.LeaveRequest
		totalCount  int
		err         error
	)

	if requesterRole == "admin" {
		leaves, totalCount, err = s.leaveRepo.FindAll(status, pageSize, offset)
	} else {
		leaves, totalCount, err = s.leaveRepo.FindByManagerID(requesterID, status, pageSize, offset)
	}
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to retrieve leave requests", err)
	}

	// Enrich with related data
	var leaveResponses []dto.LeaveRequestResponse
	for _, leave := range leaves {
		enrichedLeave, err := s.enrichLeaveRequest(leave)
		if err != nil {
			leaveResponses = append(leaveResponses, *s.mapToLeaveRequestResponse(leave))
			continue
		}
		leaveResponses = append(leaveResponses, *enrichedLeave)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &dto.LeaveRequestsListResponse{
		LeaveRequests: leaveResponses,
		TotalCount:    totalCount,
		Page:          page,
		PageSize:      pageSize,
		TotalPages:    totalPages,
	}, nil
}

// ApproveLeaveRequest approves a leave request
func (s *LeaveService) ApproveLeaveRequest(id, reviewerID int, reviewerRole string, reviewNotes *string) error {
	leave, err := s.leaveRepo.FindByID(id)
	if err != nil {
		return apperrors.NotFound("Leave request not found")
	}

	canManage, err := s.canManageLeaveRequest(leave.UserID, reviewerID, reviewerRole)
	if err != nil {
		return err
	}
	if !canManage {
		return apperrors.Forbidden("You can only review leave requests for employees assigned to you")
	}

	if !leave.IsPending() {
		return apperrors.BadRequest("Only pending leave requests can be approved")
	}

	// Update leave balance
	currentYear := utils.GetCurrentYear()
	balance, err := s.balanceRepo.FindByUserAndType(leave.UserID, leave.LeaveTypeID, currentYear)
	if err != nil {
		return apperrors.InternalServerError("Failed to retrieve leave balance", err)
	}

	// Increment used days
	if err := s.balanceRepo.IncrementUsedDays(balance.ID, leave.TotalDays); err != nil {
		return apperrors.InternalServerError("Failed to update leave balance", err)
	}

	// Update leave request status
	if err := s.leaveRepo.UpdateStatus(id, "approved", reviewerID, reviewNotes); err != nil {
		// Rollback balance update
		s.balanceRepo.DecrementUsedDays(balance.ID, leave.TotalDays)
		return apperrors.InternalServerError("Failed to approve leave request", err)
	}

	return nil
}

// RejectLeaveRequest rejects a leave request
func (s *LeaveService) RejectLeaveRequest(id, reviewerID int, reviewerRole string, reviewNotes *string) error {
	leave, err := s.leaveRepo.FindByID(id)
	if err != nil {
		return apperrors.NotFound("Leave request not found")
	}

	canManage, err := s.canManageLeaveRequest(leave.UserID, reviewerID, reviewerRole)
	if err != nil {
		return err
	}
	if !canManage {
		return apperrors.Forbidden("You can only review leave requests for employees assigned to you")
	}

	if !leave.IsPending() {
		return apperrors.BadRequest("Only pending leave requests can be rejected")
	}

	if err := s.leaveRepo.UpdateStatus(id, "rejected", reviewerID, reviewNotes); err != nil {
		return apperrors.InternalServerError("Failed to reject leave request", err)
	}

	return nil
}

func (s *LeaveService) canManageLeaveRequest(employeeID, reviewerID int, reviewerRole string) (bool, error) {
	if reviewerRole == "admin" {
		return true, nil
	}

	if reviewerRole != "manager" {
		return false, nil
	}

	employee, err := s.userRepo.FindByID(employeeID)
	if err != nil {
		return false, apperrors.NotFound("Employee not found")
	}

	if employee.ManagerID == nil {
		return false, nil
	}

	return *employee.ManagerID == reviewerID, nil
}

// GetLeaveTypes retrieves all active leave types
func (s *LeaveService) GetLeaveTypes() ([]dto.LeaveTypeResponse, error) {
	leaveTypes, err := s.leaveTypeRepo.FindAll()
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to retrieve leave types", err)
	}

	var response []dto.LeaveTypeResponse
	for _, lt := range leaveTypes {
		response = append(response, dto.LeaveTypeResponse{
			ID:          lt.ID,
			Name:        lt.Name,
			Description: lt.Description,
			IsActive:    lt.IsActive,
		})
	}

	return response, nil
}

// GetUserLeaveBalance retrieves leave balances for a user
func (s *LeaveService) GetUserLeaveBalance(userID int) ([]dto.LeaveBalanceResponse, error) {
	currentYear := utils.GetCurrentYear()
	if err := s.ensureUserHasCurrentYearBalances(userID, currentYear); err != nil {
		return nil, err
	}

	balances, err := s.balanceRepo.FindByUserID(userID, currentYear)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to retrieve leave balances", err)
	}

	var response []dto.LeaveBalanceResponse
	for _, balance := range balances {
		// Get leave type details
		leaveType, err := s.leaveTypeRepo.FindByID(balance.LeaveTypeID)
		if err != nil {
			continue
		}

		response = append(response, dto.LeaveBalanceResponse{
			ID:          balance.ID,
			UserID:      balance.UserID,
			LeaveTypeID: balance.LeaveTypeID,
			LeaveType: &dto.LeaveTypeResponse{
				ID:          leaveType.ID,
				Name:        leaveType.Name,
				Description: leaveType.Description,
				IsActive:    leaveType.IsActive,
			},
			TotalDays:     balance.TotalDays,
			UsedDays:      balance.UsedDays,
			AvailableDays: balance.AvailableDays,
			Year:          balance.Year,
		})
	}

	return response, nil
}

func (s *LeaveService) ensureUserHasCurrentYearBalances(userID, year int) error {
	leaveTypes, err := s.leaveTypeRepo.FindAll()
	if err != nil {
		return apperrors.InternalServerError("Failed to retrieve leave types", err)
	}

	for _, leaveType := range leaveTypes {
		exists, err := s.balanceRepo.BalanceExists(userID, leaveType.ID, year)
		if err != nil {
			return apperrors.InternalServerError("Failed to check leave balance", err)
		}
		if exists {
			continue
		}

		balance := &models.LeaveBalance{
			UserID:      userID,
			LeaveTypeID: leaveType.ID,
			TotalDays:   defaultLeaveDaysByTypeName(leaveType.Name),
			UsedDays:    0,
			Year:        year,
		}

		if err := s.balanceRepo.Create(balance); err != nil {
			return apperrors.InternalServerError("Failed to initialize leave balances", err)
		}
	}

	return nil
}

// enrichLeaveRequest adds user and leave type information to a leave request
func (s *LeaveService) enrichLeaveRequest(leave *models.LeaveRequest) (*dto.LeaveRequestResponse, error) {
	response := s.mapToLeaveRequestResponse(leave)

	// Get user info
	user, err := s.userRepo.FindByID(leave.UserID)
	if err == nil {
		response.User = &dto.UserBasicInfo{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			FullName:  user.GetFullName(),
			Email:     user.Email,
		}
	}

	// Get leave type info
	leaveType, err := s.leaveTypeRepo.FindByID(leave.LeaveTypeID)
	if err == nil {
		response.LeaveType = &dto.LeaveTypeResponse{
			ID:          leaveType.ID,
			Name:        leaveType.Name,
			Description: leaveType.Description,
			IsActive:    leaveType.IsActive,
		}
	}

	// Get reviewer info if reviewed
	if leave.ReviewedBy != nil {
		reviewer, err := s.userRepo.FindByID(*leave.ReviewedBy)
		if err == nil {
			response.Reviewer = &dto.UserBasicInfo{
				ID:        reviewer.ID,
				FirstName: reviewer.FirstName,
				LastName:  reviewer.LastName,
				FullName:  reviewer.GetFullName(),
				Email:     reviewer.Email,
			}
		}
	}

	return response, nil
}

// mapToLeaveRequestResponse converts a leave request model to response DTO
func (s *LeaveService) mapToLeaveRequestResponse(leave *models.LeaveRequest) *dto.LeaveRequestResponse {
	response := &dto.LeaveRequestResponse{
		ID:          leave.ID,
		UserID:      leave.UserID,
		LeaveTypeID: leave.LeaveTypeID,
		StartDate:   utils.FormatDateString(leave.StartDate),
		EndDate:     utils.FormatDateString(leave.EndDate),
		TotalDays:   leave.TotalDays,
		Reason:      leave.Reason,
		Status:      leave.Status,
		ReviewedBy:  leave.ReviewedBy,
		ReviewNotes: leave.ReviewNotes,
		CreatedAt:   leave.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   leave.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if leave.ReviewedAt != nil {
		reviewedAt := leave.ReviewedAt.Format("2006-01-02 15:04:05")
		response.ReviewedAt = &reviewedAt
	}

	return response
}
