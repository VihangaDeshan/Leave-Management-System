package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"leave-management-system/internal/models"
	"time"
)

type LeaveRepository struct {
	db *sql.DB
}

func NewLeaveRepository(db *sql.DB) *LeaveRepository {
	return &LeaveRepository{db: db}
}

// Create creates a new leave request
func (r *LeaveRepository) Create(leave *models.LeaveRequest) error {
	query := `
		INSERT INTO leave_requests (user_id, leave_type_id, start_date, end_date, total_days, reason, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(
		query,
		leave.UserID,
		leave.LeaveTypeID,
		leave.StartDate,
		leave.EndDate,
		leave.TotalDays,
		leave.Reason,
		leave.Status,
	).Scan(&leave.ID, &leave.CreatedAt, &leave.UpdatedAt)

	return err
}

// FindByID finds a leave request by ID
func (r *LeaveRepository) FindByID(id int) (*models.LeaveRequest, error) {
	leave := &models.LeaveRequest{}
	query := `
		SELECT id, user_id, leave_type_id, start_date, end_date, total_days, reason, status,
		       reviewed_by, reviewed_at, review_notes, created_at, updated_at
		FROM leave_requests
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&leave.ID,
		&leave.UserID,
		&leave.LeaveTypeID,
		&leave.StartDate,
		&leave.EndDate,
		&leave.TotalDays,
		&leave.Reason,
		&leave.Status,
		&leave.ReviewedBy,
		&leave.ReviewedAt,
		&leave.ReviewNotes,
		&leave.CreatedAt,
		&leave.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("leave request not found")
	}
	if err != nil {
		return nil, err
	}

	return leave, nil
}

// FindByUserID finds all leave requests by user ID with pagination
func (r *LeaveRepository) FindByUserID(userID, limit, offset int) ([]*models.LeaveRequest, int, error) {
	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM leave_requests WHERE user_id = $1`
	err := r.db.QueryRow(countQuery, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get leave requests
	query := `
		SELECT id, user_id, leave_type_id, start_date, end_date, total_days, reason, status,
		       reviewed_by, reviewed_at, review_notes, created_at, updated_at
		FROM leave_requests
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var leaves []*models.LeaveRequest
	for rows.Next() {
		leave := &models.LeaveRequest{}
		err := rows.Scan(
			&leave.ID,
			&leave.UserID,
			&leave.LeaveTypeID,
			&leave.StartDate,
			&leave.EndDate,
			&leave.TotalDays,
			&leave.Reason,
			&leave.Status,
			&leave.ReviewedBy,
			&leave.ReviewedAt,
			&leave.ReviewNotes,
			&leave.CreatedAt,
			&leave.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		leaves = append(leaves, leave)
	}

	return leaves, totalCount, nil
}

// FindAll finds all leave requests with optional filtering
func (r *LeaveRepository) FindAll(status string, limit, offset int) ([]*models.LeaveRequest, int, error) {
	// Build query with optional status filter
	countQuery := `SELECT COUNT(*) FROM leave_requests`
	query := `
		SELECT id, user_id, leave_type_id, start_date, end_date, total_days, reason, status,
		       reviewed_by, reviewed_at, review_notes, created_at, updated_at
		FROM leave_requests
	`

	var args []interface{}
	argIndex := 1

	if status != "" {
		countQuery += fmt.Sprintf(" WHERE status = $%d", argIndex)
		query += fmt.Sprintf(" WHERE status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	// Get total count
	var totalCount int
	var err error
	if len(args) > 0 {
		err = r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	} else {
		err = r.db.QueryRow(countQuery).Scan(&totalCount)
	}
	if err != nil {
		return nil, 0, err
	}

	// Add ordering and pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	// Execute query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var leaves []*models.LeaveRequest
	for rows.Next() {
		leave := &models.LeaveRequest{}
		err := rows.Scan(
			&leave.ID,
			&leave.UserID,
			&leave.LeaveTypeID,
			&leave.StartDate,
			&leave.EndDate,
			&leave.TotalDays,
			&leave.Reason,
			&leave.Status,
			&leave.ReviewedBy,
			&leave.ReviewedAt,
			&leave.ReviewNotes,
			&leave.CreatedAt,
			&leave.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		leaves = append(leaves, leave)
	}

	return leaves, totalCount, nil
}

// Update updates a leave request
func (r *LeaveRepository) Update(leave *models.LeaveRequest) error {
	query := `
		UPDATE leave_requests
		SET leave_type_id = $1, start_date = $2, end_date = $3, total_days = $4, reason = $5,
		    status = $6, reviewed_by = $7, reviewed_at = $8, review_notes = $9
		WHERE id = $10
		RETURNING updated_at
	`
	return r.db.QueryRow(
		query,
		leave.LeaveTypeID,
		leave.StartDate,
		leave.EndDate,
		leave.TotalDays,
		leave.Reason,
		leave.Status,
		leave.ReviewedBy,
		leave.ReviewedAt,
		leave.ReviewNotes,
		leave.ID,
	).Scan(&leave.UpdatedAt)
}

// UpdateStatus updates only the status of a leave request
func (r *LeaveRepository) UpdateStatus(id int, status string, reviewedBy int, reviewNotes *string) error {
	query := `
		UPDATE leave_requests
		SET status = $1, reviewed_by = $2, reviewed_at = $3, review_notes = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(query, status, reviewedBy, time.Now(), reviewNotes, id)
	return err
}

// Delete deletes a leave request (soft delete by changing status to cancelled)
func (r *LeaveRepository) Delete(id int) error {
	query := `UPDATE leave_requests SET status = 'cancelled' WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// CheckOverlap checks if there's an overlapping leave request for a user
func (r *LeaveRepository) CheckOverlap(userID int, startDate, endDate time.Time, excludeID *int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM leave_requests
			WHERE user_id = $1
			AND status != 'rejected' AND status != 'cancelled'
			AND (
				(start_date <= $2 AND end_date >= $2) OR
				(start_date <= $3 AND end_date >= $3) OR
				(start_date >= $2 AND end_date <= $3)
			)
	`

	args := []interface{}{userID, startDate, endDate}

	if excludeID != nil {
		query += ` AND id != $4`
		args = append(args, *excludeID)
	}

	query += `)`

	var exists bool
	err := r.db.QueryRow(query, args...).Scan(&exists)
	return exists, err
}

// GetPendingCount returns the count of pending leave requests
func (r *LeaveRepository) GetPendingCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM leave_requests WHERE status = 'pending'`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// GetEmployeesOnLeave returns users currently on leave
func (r *LeaveRepository) GetEmployeesOnLeave(date time.Time) ([]*models.LeaveRequest, error) {
	query := `
		SELECT id, user_id, leave_type_id, start_date, end_date, total_days, reason, status,
		       reviewed_by, reviewed_at, review_notes, created_at, updated_at
		FROM leave_requests
		WHERE status = 'approved'
		AND start_date <= $1
		AND end_date >= $1
		ORDER BY start_date
	`
	rows, err := r.db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaves []*models.LeaveRequest
	for rows.Next() {
		leave := &models.LeaveRequest{}
		err := rows.Scan(
			&leave.ID,
			&leave.UserID,
			&leave.LeaveTypeID,
			&leave.StartDate,
			&leave.EndDate,
			&leave.TotalDays,
			&leave.Reason,
			&leave.Status,
			&leave.ReviewedBy,
			&leave.ReviewedAt,
			&leave.ReviewNotes,
			&leave.CreatedAt,
			&leave.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		leaves = append(leaves, leave)
	}

	return leaves, nil
}
