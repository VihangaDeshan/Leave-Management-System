-- Leave Management System Database Schema
-- Version: 1.0.0
-- Description: Initial database setup with all tables, constraints, and indexes

-- ============================================
-- Table: users
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('employee', 'manager', 'admin')),
    department VARCHAR(100),
    manager_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for users table
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_manager_id ON users(manager_id);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);

-- ============================================
-- Table: leave_types
-- ============================================
CREATE TABLE IF NOT EXISTS leave_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for leave_types table
CREATE INDEX IF NOT EXISTS idx_leave_types_is_active ON leave_types(is_active);

-- ============================================
-- Table: leave_balances
-- ============================================
CREATE TABLE IF NOT EXISTS leave_balances (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    leave_type_id INTEGER NOT NULL REFERENCES leave_types(id) ON DELETE CASCADE,
    total_days DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    used_days DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    year INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_leave_type_year UNIQUE (user_id, leave_type_id, year),
    CONSTRAINT chk_positive_days CHECK (total_days >= 0 AND used_days >= 0),
    CONSTRAINT chk_used_not_exceed_total CHECK (used_days <= total_days)
);

-- Indexes for leave_balances table
CREATE INDEX IF NOT EXISTS idx_leave_balances_user_id ON leave_balances(user_id);
CREATE INDEX IF NOT EXISTS idx_leave_balances_leave_type_id ON leave_balances(leave_type_id);
CREATE INDEX IF NOT EXISTS idx_leave_balances_year ON leave_balances(year);
CREATE INDEX IF NOT EXISTS idx_leave_balances_user_year ON leave_balances(user_id, year);

-- ============================================
-- Table: leave_requests
-- ============================================
CREATE TABLE IF NOT EXISTS leave_requests (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    leave_type_id INTEGER NOT NULL REFERENCES leave_types(id) ON DELETE RESTRICT,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_days DECIMAL(5,2) NOT NULL,
    reason TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled')),
    reviewed_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMP WITH TIME ZONE,
    review_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_date_order CHECK (end_date >= start_date),
    CONSTRAINT chk_positive_total_days CHECK (total_days > 0)
);

-- Indexes for leave_requests table
CREATE INDEX IF NOT EXISTS idx_leave_requests_user_id ON leave_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_leave_requests_leave_type_id ON leave_requests(leave_type_id);
CREATE INDEX IF NOT EXISTS idx_leave_requests_status ON leave_requests(status);
CREATE INDEX IF NOT EXISTS idx_leave_requests_start_date ON leave_requests(start_date);
CREATE INDEX IF NOT EXISTS idx_leave_requests_end_date ON leave_requests(end_date);
CREATE INDEX IF NOT EXISTS idx_leave_requests_reviewed_by ON leave_requests(reviewed_by);
CREATE INDEX IF NOT EXISTS idx_leave_requests_user_status ON leave_requests(user_id, status);

-- ============================================
-- Table: audit_logs
-- ============================================
CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL CHECK (action IN ('create', 'update', 'delete', 'approve', 'reject', 'cancel')),
    entity_type VARCHAR(50) NOT NULL,
    entity_id INTEGER NOT NULL,
    old_value JSONB,
    new_value JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for audit_logs table
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);

-- ============================================
-- Trigger: Update updated_at timestamp
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to relevant tables
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_leave_balances_updated_at
    BEFORE UPDATE ON leave_balances
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_leave_requests_updated_at
    BEFORE UPDATE ON leave_requests
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- Seed Data: Insert default leave types
-- ============================================
INSERT INTO leave_types (name, description, is_active) VALUES
    ('Annual Leave', 'Regular paid vacation leave', true),
    ('Sick Leave', 'Leave for medical reasons', true),
    ('Casual Leave', 'Short-term unplanned leave', true),
    ('Maternity Leave', 'Leave for maternity purposes', true),
    ('Paternity Leave', 'Leave for paternity purposes', true),
    ('Unpaid Leave', 'Leave without pay', true)
ON CONFLICT (name) DO NOTHING;

-- ============================================
-- Seed Data: Insert default admin user
-- Password: admin123 (hashed with bcrypt)
-- IMPORTANT: Change this password in production!
-- ============================================
INSERT INTO users (email, password_hash, first_name, last_name, role, is_active) VALUES
    ('admin@example.com', '$2a$10$6BqCVZqPOz8EKSGqDDvvx.LmfKj2uTJXxLhHCPuMYf3n9kE9gRfFO', 'System', 'Admin', 'admin', true)
ON CONFLICT (email) DO NOTHING;

-- ============================================
-- Success Message
-- ============================================
DO $$
BEGIN
    RAISE NOTICE 'Database schema created successfully!';
    RAISE NOTICE 'Default admin credentials:';
    RAISE NOTICE '  Email: admin@example.com';
    RAISE NOTICE '  Password: admin123';
    RAISE NOTICE 'IMPORTANT: Change the admin password immediately in production!';
END $$;
