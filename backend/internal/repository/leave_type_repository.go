package repository

import (
	"database/sql"
	"errors"
	"leave-management-system/internal/models"
)

type LeaveTypeRepository struct {
	db *sql.DB
}

func NewLeaveTypeRepository(db *sql.DB) *LeaveTypeRepository {
	return &LeaveTypeRepository{db: db}
}

// FindAll returns all leave types
func (r *LeaveTypeRepository) FindAll() ([]*models.LeaveType, error) {
	query := `
		SELECT id, name, description, is_active, created_at
		FROM leave_types
		WHERE is_active = true
		ORDER BY name
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaveTypes []*models.LeaveType
	for rows.Next() {
		lt := &models.LeaveType{}
		err := rows.Scan(
			&lt.ID,
			&lt.Name,
			&lt.Description,
			&lt.IsActive,
			&lt.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		leaveTypes = append(leaveTypes, lt)
	}

	return leaveTypes, nil
}

// FindByID finds a leave type by ID
func (r *LeaveTypeRepository) FindByID(id int) (*models.LeaveType, error) {
	lt := &models.LeaveType{}
	query := `
		SELECT id, name, description, is_active, created_at
		FROM leave_types
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&lt.ID,
		&lt.Name,
		&lt.Description,
		&lt.IsActive,
		&lt.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("leave type not found")
	}
	if err != nil {
		return nil, err
	}

	return lt, nil
}

// Create creates a new leave type
func (r *LeaveTypeRepository) Create(leaveType *models.LeaveType) error {
	query := `
		INSERT INTO leave_types (name, description, is_active)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(
		query,
		leaveType.Name,
		leaveType.Description,
		leaveType.IsActive,
	).Scan(&leaveType.ID, &leaveType.CreatedAt)

	return err
}

// Update updates a leave type
func (r *LeaveTypeRepository) Update(leaveType *models.LeaveType) error {
	query := `
		UPDATE leave_types
		SET name = $1, description = $2, is_active = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(
		query,
		leaveType.Name,
		leaveType.Description,
		leaveType.IsActive,
		leaveType.ID,
	)
	return err
}

// Delete soft deletes a leave type
func (r *LeaveTypeRepository) Delete(id int) error {
	query := `UPDATE leave_types SET is_active = false WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
