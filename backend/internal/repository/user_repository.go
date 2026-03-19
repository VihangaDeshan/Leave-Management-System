package repository

import (
	"database/sql"
	"errors"
	"leave-management-system/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, role, department, manager_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(
		query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Role,
		user.Department,
		user.ManagerID,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, department, manager_id, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Department,
		&user.ManagerID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, department, manager_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Department,
		&user.ManagerID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindAll finds all users with pagination
func (r *UserRepository) FindAll(limit, offset int) ([]*models.User, int, error) {
	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM users`
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get users
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, department, manager_id, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.Department,
			&user.ManagerID,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, totalCount, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, department = $3, manager_id = $4, role = $5, is_active = $6
		WHERE id = $7
		RETURNING updated_at
	`
	return r.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Department,
		user.ManagerID,
		user.Role,
		user.IsActive,
		user.ID,
	).Scan(&user.UpdatedAt)
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(userID int, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1 WHERE id = $2`
	_, err := r.db.Exec(query, passwordHash, userID)
	return err
}

// Delete soft deletes a user by setting is_active to false
func (r *UserRepository) Delete(id int) error {
	query := `UPDATE users SET is_active = false WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// EmailExists checks if an email already exists
func (r *UserRepository) EmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

// GetActiveEmployeeCount returns the count of active employees
func (r *UserRepository) GetActiveEmployeeCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE is_active = true`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// GetDistinctDepartments returns a list of unique departments
func (r *UserRepository) GetDistinctDepartments() ([]string, error) {
	query := `SELECT DISTINCT department FROM users WHERE department IS NOT NULL AND department != '' ORDER BY department`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []string
	for rows.Next() {
		var dept string
		if err := rows.Scan(&dept); err != nil {
			return nil, err
		}
		departments = append(departments, dept)
	}
	return departments, nil
}
