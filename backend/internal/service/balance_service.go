package service

import (
	"leave-management-system/internal/dto"
	"leave-management-system/internal/models"
	"leave-management-system/internal/repository"
	apperrors "leave-management-system/pkg/errors"
)

type BalanceService struct {
	balanceRepo   *repository.BalanceRepository
	leaveTypeRepo *repository.LeaveTypeRepository
	userRepo      *repository.UserRepository
}

func NewBalanceService(
	balanceRepo *repository.BalanceRepository,
	leaveTypeRepo *repository.LeaveTypeRepository,
	userRepo *repository.UserRepository,
) *BalanceService {
	return &BalanceService{
		balanceRepo:   balanceRepo,
		leaveTypeRepo: leaveTypeRepo,
		userRepo:      userRepo,
	}
}

// AllocateBalance allocates leave balance for a user
func (s *BalanceService) AllocateBalance(req *dto.AllocateBalanceRequest) (*dto.LeaveBalanceResponse, error) {
	// Verify user exists
	_, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, apperrors.NotFound("User not found")
	}

	// Verify leave type exists
	leaveType, err := s.leaveTypeRepo.FindByID(req.LeaveTypeID)
	if err != nil {
		return nil, apperrors.NotFound("Leave type not found")
	}

	// Check if balance already exists
	exists, err := s.balanceRepo.BalanceExists(req.UserID, req.LeaveTypeID, req.Year)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to check balance existence", err)
	}
	if exists {
		return nil, apperrors.Conflict("Balance for this leave type and year already exists")
	}

	// Create balance
	balance := &models.LeaveBalance{
		UserID:      req.UserID,
		LeaveTypeID: req.LeaveTypeID,
		TotalDays:   req.TotalDays,
		UsedDays:    0,
		Year:        req.Year,
	}

	if err := s.balanceRepo.Create(balance); err != nil {
		return nil, apperrors.InternalServerError("Failed to allocate balance", err)
	}

	return &dto.LeaveBalanceResponse{
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
	}, nil
}

// UpdateBalance updates a user's leave balance
func (s *BalanceService) UpdateBalance(balanceID int, req *dto.UpdateBalanceRequest) (*dto.LeaveBalanceResponse, error) {
	// Find existing balance
	balance, err := s.balanceRepo.FindByUserAndType(0, 0, 0) // We'll need to query by ID
	if err != nil {
		return nil, apperrors.NotFound("Balance not found")
	}

	// Update total days
	balance.TotalDays = req.TotalDays

	if err := s.balanceRepo.Update(balance); err != nil {
		return nil, apperrors.InternalServerError("Failed to update balance", err)
	}

	// Get leave type for response
	leaveType, _ := s.leaveTypeRepo.FindByID(balance.LeaveTypeID)

	response := &dto.LeaveBalanceResponse{
		ID:            balance.ID,
		UserID:        balance.UserID,
		LeaveTypeID:   balance.LeaveTypeID,
		TotalDays:     balance.TotalDays,
		UsedDays:      balance.UsedDays,
		AvailableDays: balance.AvailableDays,
		Year:          balance.Year,
	}

	if leaveType != nil {
		response.LeaveType = &dto.LeaveTypeResponse{
			ID:          leaveType.ID,
			Name:        leaveType.Name,
			Description: leaveType.Description,
			IsActive:    leaveType.IsActive,
		}
	}

	return response, nil
}
