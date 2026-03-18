package handler

import (
	"encoding/json"
	"leave-management-system/internal/dto"
	"leave-management-system/internal/middleware"
	"leave-management-system/internal/service"
	"leave-management-system/internal/utils"
	apperrors "leave-management-system/pkg/errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService *service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validate.Struct(&req); err != nil {
		utils.ErrorResponse(w, "Validation error", err.Error(), http.StatusBadRequest)
		return
	}

	// Call service
	response, err := h.authService.Register(&req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Registration successful", response, http.StatusCreated)
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validate.Struct(&req); err != nil {
		utils.ErrorResponse(w, "Validation error", err.Error(), http.StatusBadRequest)
		return
	}

	// Call service
	response, err := h.authService.Login(&req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Login successful", response, http.StatusOK)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validate.Struct(&req); err != nil {
		utils.ErrorResponse(w, "Validation error", err.Error(), http.StatusBadRequest)
		return
	}

	// Call service
	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Token refreshed successfully", response, http.StatusOK)
}

// GetMe returns current user information
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.ErrorResponse(w, "Unauthorized", nil, http.StatusUnauthorized)
		return
	}

	// Call service
	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "User retrieved successfully", user, http.StatusOK)
}
