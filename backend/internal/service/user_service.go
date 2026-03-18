package service

import (
	"leave-management-system/internal/dto"
	"leave-management-system/internal/models"
	"leave-management-system/internal/repository"
	"leave-management-system/internal/utils"
	apperrors "leave-management-system/pkg/errors"
	"math"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetProfile retrieves a user's profile
func (s *UserService) GetProfile(userID int) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperrors.NotFound("User not found")
	}

	return s.mapToUserResponse(user), nil
}

// UpdateProfile updates a user's own profile
func (s *UserService) UpdateProfile(userID int, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperrors.NotFound("User not found")
	}

	// Update fields
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Department = req.Department

	if err := s.userRepo.Update(user); err != nil {
		return nil, apperrors.InternalServerError("Failed to update profile", err)
	}

	return s.mapToUserResponse(user), nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(userID int, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return apperrors.NotFound("User not found")
	}

	// Verify old password
	if err := utils.ComparePassword(user.PasswordHash, req.OldPassword); err != nil {
		return apperrors.BadRequest("Current password is incorrect")
	}

	// Hash new password
	newPasswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return apperrors.InternalServerError("Failed to hash password", err)
	}

	// Update password
	if err := s.userRepo.UpdatePassword(userID, newPasswordHash); err != nil {
		return apperrors.InternalServerError("Failed to update password", err)
	}

	return nil
}

// GetAllUsers retrieves all users with pagination
func (s *UserService) GetAllUsers(page, pageSize int) (*dto.UsersListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	users, totalCount, err := s.userRepo.FindAll(pageSize, offset)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to retrieve users", err)
	}

	// Map to response
	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, *s.mapToUserResponse(user))
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &dto.UsersListResponse{
		Users:      userResponses,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id int) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFound("User not found")
	}

	return s.mapToUserResponse(user), nil
}

// CreateUser creates a new user (admin only)
func (s *UserService) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if email already exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to check email existence", err)
	}
	if exists {
		return nil, apperrors.Conflict("Email already registered")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to hash password", err)
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
		Department:   req.Department,
		ManagerID:    req.ManagerID,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, apperrors.InternalServerError("Failed to create user", err)
	}

	return s.mapToUserResponse(user), nil
}

// UpdateUser updates a user (admin only)
func (s *UserService) UpdateUser(id int, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFound("User not found")
	}

	// Update fields if provided
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Department != nil {
		user.Department = req.Department
	}
	if req.ManagerID != nil {
		user.ManagerID = req.ManagerID
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, apperrors.InternalServerError("Failed to update user", err)
	}

	return s.mapToUserResponse(user), nil
}

// DeactivateUser deactivates a user (soft delete)
func (s *UserService) DeactivateUser(id int) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return apperrors.NotFound("User not found")
	}

	if !user.IsActive {
		return apperrors.BadRequest("User is already deactivated")
	}

	if err := s.userRepo.Delete(id); err != nil {
		return apperrors.InternalServerError("Failed to deactivate user", err)
	}

	return nil
}

// Helper function to map User model to UserResponse DTO
func (s *UserService) mapToUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		FullName:   user.GetFullName(),
		Role:       user.Role,
		Department: user.Department,
		ManagerID:  user.ManagerID,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
