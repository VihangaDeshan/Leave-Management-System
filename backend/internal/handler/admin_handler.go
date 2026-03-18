package handler

import (
	"encoding/json"
	"leave-management-system/internal/dto"
	"leave-management-system/internal/middleware"
	"leave-management-system/internal/service"
	"leave-management-system/internal/utils"
	apperrors "leave-management-system/pkg/errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type AdminHandler struct {
	leaveService     *service.LeaveService
	userService      *service.UserService
	balanceService   *service.BalanceService
	dashboardService *service.DashboardService
	validate         *validator.Validate
}

func NewAdminHandler(
	leaveService *service.LeaveService,
	userService *service.UserService,
	balanceService *service.BalanceService,
	dashboardService *service.DashboardService,
) *AdminHandler {
	return &AdminHandler{
		leaveService:     leaveService,
		userService:      userService,
		balanceService:   balanceService,
		dashboardService: dashboardService,
		validate:         validator.New(),
	}
}

// GetAllLeaves retrieves all leave requests
func (h *AdminHandler) GetAllLeaves(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	status := r.URL.Query().Get("status")

	response, err := h.leaveService.GetAllLeaveRequests(status, page, pageSize)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Leave requests retrieved successfully", response, http.StatusOK)
}

// ApproveLeave approves a leave request
func (h *AdminHandler) ApproveLeave(w http.ResponseWriter, r *http.Request) {
	reviewerID, _ := middleware.GetUserIDFromContext(r.Context())

	vars := mux.Vars(r)
	leaveID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, "Invalid leave ID", nil, http.StatusBadRequest)
		return
	}

	var req dto.ReviewLeaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.ReviewNotes != nil {
		err = h.leaveService.ApproveLeaveRequest(leaveID, reviewerID, req.ReviewNotes)
	} else {
		err = h.leaveService.ApproveLeaveRequest(leaveID, reviewerID, nil)
	}

	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Leave request approved successfully", nil, http.StatusOK)
}

// RejectLeave rejects a leave request
func (h *AdminHandler) RejectLeave(w http.ResponseWriter, r *http.Request) {
	reviewerID, _ := middleware.GetUserIDFromContext(r.Context())

	vars := mux.Vars(r)
	leaveID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, "Invalid leave ID", nil, http.StatusBadRequest)
		return
	}

	var req dto.ReviewLeaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.ReviewNotes != nil {
		err = h.leaveService.RejectLeaveRequest(leaveID, reviewerID, req.ReviewNotes)
	} else {
		err = h.leaveService.RejectLeaveRequest(leaveID, reviewerID, nil)
	}

	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Leave request rejected successfully", nil, http.StatusOK)
}

// GetAllUsers retrieves all users
func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	response, err := h.userService.GetAllUsers(page, pageSize)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Users retrieved successfully", response, http.StatusOK)
}

// CreateUser creates a new user
func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		utils.ErrorResponse(w, "Validation error", err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.userService.CreateUser(&req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "User created successfully", response, http.StatusCreated)
}

// UpdateUser updates a user
func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, "Invalid user ID", nil, http.StatusBadRequest)
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	response, err := h.userService.UpdateUser(userID, &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "User updated successfully", response, http.StatusOK)
}

// Get Dashboard Stats
func (h *AdminHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	response, err := h.dashboardService.GetDashboardStats()
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Dashboard stats retrieved successfully", response, http.StatusOK)
}

// AllocateBalance allocates leave balance to a user
func (h *AdminHandler) AllocateBalance(w http.ResponseWriter, r *http.Request) {
	var req dto.AllocateBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		utils.ErrorResponse(w, "Validation error", err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.balanceService.AllocateBalance(&req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Balance allocated successfully", response, http.StatusCreated)
}
