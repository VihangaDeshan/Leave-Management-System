package repository

import (
	"database/sql"
	"errors"
	"leave-management-system/internal/models"
)

type BalanceRepository struct {
	db *sql.DB
}

func NewBalanceRepository(db *sql.DB) *BalanceRepository {
	return &BalanceRepository{db: db}
}

// Create creates a new leave balance entry
func (r *BalanceRepository) Create(balance *models.LeaveBalance) error {
	query := `
		INSERT INTO leave_balances (user_id, leave_type_id, total_days, used_days, year)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, available_days, created_at, updated_at
	`
	err := r.db.QueryRow(
		query,
		balance.UserID,
		balance.LeaveTypeID,
		balance.TotalDays,
		balance.UsedDays,
		balance.Year,
	).Scan(&balance.ID, &balance.AvailableDays, &balance.CreatedAt, &balance.UpdatedAt)

	return err
}

// FindByUserAndType finds a leave balance by user ID, leave type ID, and year
func (r *BalanceRepository) FindByUserAndType(userID, leaveTypeID, year int) (*models.LeaveBalance, error) {
	balance := &models.LeaveBalance{}
	query := `
		SELECT id, user_id, leave_type_id, total_days, used_days, available_days, year, created_at, updated_at
		FROM leave_balances
		WHERE user_id = $1 AND leave_type_id = $2 AND year = $3
	`
	err := r.db.QueryRow(query, userID, leaveTypeID, year).Scan(
		&balance.ID,
		&balance.UserID,
		&balance.LeaveTypeID,
		&balance.TotalDays,
		&balance.UsedDays,
		&balance.AvailableDays,
		&balance.Year,
		&balance.CreatedAt,
		&balance.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("balance not found")
	}
	if err != nil {
		return nil, err
	}

	return balance, nil
}

// FindByUserID finds all leave balances for a user in a given year
func (r *BalanceRepository) FindByUserID(userID, year int) ([]*models.LeaveBalance, error) {
	query := `
		SELECT id, user_id, leave_type_id, total_days, used_days, available_days, year, created_at, updated_at
		FROM leave_balances
		WHERE user_id = $1 AND year = $2
		ORDER BY leave_type_id
	`
	rows, err := r.db.Query(query, userID, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []*models.LeaveBalance
	for rows.Next() {
		balance := &models.LeaveBalance{}
		err := rows.Scan(
			&balance.ID,
			&balance.UserID,
			&balance.LeaveTypeID,
			&balance.TotalDays,
			&balance.UsedDays,
			&balance.AvailableDays,
			&balance.Year,
			&balance.CreatedAt,
			&balance.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}

	return balances, nil
}

// Update updates a leave balance
func (r *BalanceRepository) Update(balance *models.LeaveBalance) error {
	query := `
		UPDATE leave_balances
		SET total_days = $1, used_days = $2
		WHERE id = $3
		RETURNING available_days, updated_at
	`
	return r.db.QueryRow(
		query,
		balance.TotalDays,
		balance.UsedDays,
		balance.ID,
	).Scan(&balance.AvailableDays, &balance.UpdatedAt)
}

// UpdateUsedDays updates only the used days of a balance
func (r *BalanceRepository) UpdateUsedDays(balanceID int, usedDays float64) error {
	query := `UPDATE leave_balances SET used_days = $1 WHERE id = $2`
	_, err := r.db.Exec(query, usedDays, balanceID)
	return err
}

// IncrementUsedDays increments the used days by a given amount
func (r *BalanceRepository) IncrementUsedDays(balanceID int, days float64) error {
	query := `UPDATE leave_balances SET used_days = used_days + $1 WHERE id = $2`
	_, err := r.db.Exec(query, days, balanceID)
	return err
}

// DecrementUsedDays decrements the used days by a given amount
func (r *BalanceRepository) DecrementUsedDays(balanceID int, days float64) error {
	query := `UPDATE leave_balances SET used_days = used_days - $1 WHERE id = $2`
	_, err := r.db.Exec(query, days, balanceID)
	return err
}

// Delete deletes a leave balance
func (r *BalanceRepository) Delete(id int) error {
	query := `DELETE FROM leave_balances WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// FindAll returns all balances with optional filtering
func (r *BalanceRepository) FindAll(year int, limit, offset int) ([]*models.LeaveBalance, int, error) {
	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM leave_balances WHERE year = $1`
	err := r.db.QueryRow(countQuery, year).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get balances
	query := `
		SELECT id, user_id, leave_type_id, total_days, used_days, available_days, year, created_at, updated_at
		FROM leave_balances
		WHERE year = $1
		ORDER BY user_id, leave_type_id
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, year, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var balances []*models.LeaveBalance
	for rows.Next() {
		balance := &models.LeaveBalance{}
		err := rows.Scan(
			&balance.ID,
			&balance.UserID,
			&balance.LeaveTypeID,
			&balance.TotalDays,
			&balance.UsedDays,
			&balance.AvailableDays,
			&balance.Year,
			&balance.CreatedAt,
			&balance.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		balances = append(balances, balance)
	}

	return balances, totalCount, nil
}

// BalanceExists checks if a balance entry exists
func (r *BalanceRepository) BalanceExists(userID, leaveTypeID, year int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM leave_balances WHERE user_id = $1 AND leave_type_id = $2 AND year = $3)`
	err := r.db.QueryRow(query, userID, leaveTypeID, year).Scan(&exists)
	return exists, err
}
