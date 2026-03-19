package service

import (
	"errors"
	"leave-management-system/internal/dto"
	"leave-management-system/internal/models"
	"leave-management-system/internal/repository"
	"leave-management-system/internal/utils"
	apperrors "leave-management-system/pkg/errors"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register registers a new user
func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.LoginResponse, error) {
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
		Role:         "employee", // Default role
		Department:   req.Department,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, apperrors.InternalServerError("Failed to create user", err)
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(user)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to generate access token", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(user)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to generate refresh token", err)
	}

	// Build response
	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserInfo{
			ID:         user.ID,
			Email:      user.Email,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Role:       user.Role,
			Department: user.Department,
			IsActive:   user.IsActive,
		},
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, apperrors.Unauthorized("Invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, apperrors.Forbidden("Account is disabled")
	}

	// Compare password
	if err := utils.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return nil, apperrors.Unauthorized("Invalid email or password")
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(user)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to generate access token", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(user)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to generate refresh token", err)
	}

	// Build response
	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserInfo{
			ID:         user.ID,
			Email:      user.Email,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Role:       user.Role,
			Department: user.Department,
			IsActive:   user.IsActive,
		},
	}, nil
}

// RefreshToken generates a new access token from a refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*dto.RefreshTokenResponse, error) {
	// Verify refresh token
	claims, err := utils.VerifyToken(refreshToken)
	if err != nil {
		return nil, apperrors.Unauthorized("Invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, apperrors.Unauthorized("User not found")
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, apperrors.Forbidden("Account is disabled")
	}

	// Generate new access token
	accessToken, err := utils.GenerateToken(user)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to generate access token", err)
	}

	return &dto.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

// GetCurrentUser retrieves the current authenticated user
func (s *AuthService) GetCurrentUser(userID int) (*dto.UserInfo, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, errors.New("user not found")) {
			return nil, apperrors.NotFound("User not found")
		}
		return nil, apperrors.InternalServerError("Failed to retrieve user", err)
	}

	return &dto.UserInfo{
		ID:         user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Role:       user.Role,
		Department: user.Department,
		IsActive:   user.IsActive,
	}, nil
}

// GetDepartments retrieves a list of all distinct departments
func (s *AuthService) GetDepartments() ([]string, error) {
	departments, err := s.userRepo.GetDistinctDepartments()
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to retrieve departments", err)
	}
	return departments, nil
}
