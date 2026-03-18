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

type LeaveHandler struct {
	leaveService *service.LeaveService
	validate     *validator.Validate
}

func NewLeaveHandler(leaveService *service.LeaveService) *LeaveHandler {
	return &LeaveHandler{
		leaveService: leaveService,
		validate:     validator.New(),
	}
}

// CreateLeave handles creating a new leave request
func (h *LeaveHandler) CreateLeave(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var req dto.CreateLeaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		utils.ErrorResponse(w, "Validation error", err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.leaveService.CreateLeaveRequest(userID, &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Leave request created successfully", response, http.StatusCreated)
}

// GetUserLeaves retrieves all leave requests for the current user
func (h *LeaveHandler) GetUserLeaves(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	response, err := h.leaveService.GetUserLeaveRequests(userID, page, pageSize)
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

// GetLeaveByID retrieves a specific leave request
func (h *LeaveHandler) GetLeaveByID(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	role, _ := middleware.GetUserRoleFromContext(r.Context())

	vars := mux.Vars(r)
	leaveID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, "Invalid leave ID", nil, http.StatusBadRequest)
		return
	}

	isAdmin := role == "admin" || role == "manager"
	response, err := h.leaveService.GetLeaveRequestByID(leaveID, userID, isAdmin)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Leave request retrieved successfully", response, http.StatusOK)
}

// CancelLeave cancels a pending leave request
func (h *LeaveHandler) CancelLeave(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	role, _ := middleware.GetUserRoleFromContext(r.Context())

	vars := mux.Vars(r)
	leaveID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, "Invalid leave ID", nil, http.StatusBadRequest)
		return
	}

	isAdmin := role == "admin" || role == "manager"
	if err := h.leaveService.CancelLeaveRequest(leaveID, userID, isAdmin); err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Leave request cancelled successfully", nil, http.StatusOK)
}

// GetLeaveTypes retrieves all available leave types
func (h *LeaveHandler) GetLeaveTypes(w http.ResponseWriter, r *http.Request) {
	response, err := h.leaveService.GetLeaveTypes()
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Leave types retrieved successfully", response, http.StatusOK)
}

// GetUserBalance retrieves leave balance for the current user
func (h *LeaveHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	response, err := h.leaveService.GetUserLeaveBalance(userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			utils.ErrorResponse(w, appErr.Message, nil, appErr.Code)
			return
		}
		utils.ErrorResponse(w, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Leave balance retrieved successfully", response, http.StatusOK)
}
