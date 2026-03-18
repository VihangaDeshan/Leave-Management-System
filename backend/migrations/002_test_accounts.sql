-- Test Accounts Seed Data
-- Description: Creates test accounts with sample leave balances for development/testing

-- ============================================
-- Insert Test Users
-- Password: admin123, manager123, employee123
-- (Hashed with bcrypt cost 10)
-- ============================================

-- Delete existing test accounts if they exist
DELETE FROM users WHERE email IN (
    'admin@abccompany.com',
    'manager@abccompany.com',
    'employee@abccompany.com'
);

-- Insert Admin Account (Password: admin123)
INSERT INTO users (email, password_hash, first_name, last_name, role, department, is_active) VALUES
    ('admin@abccompany.com', '$2a$10$F9udk7sFvFBN8HehqCwDKOaX.XL7TWLtM/AX3AchKeg.owQ7Bi7LC', 'Admin', 'User', 'admin', 'Administration', true);

-- Insert Manager Account (Password: manager123)
INSERT INTO users (email, password_hash, first_name, last_name, role, department, is_active) VALUES
    ('manager@abccompany.com', '$2a$10$W5cFn6eO1HWm1UzPNKyxie.WJg/tTp7PxBfhSX2CEMnS97GPFoJcS', 'Manager', 'User', 'manager', 'Human Resources', true);

-- Insert Employee Account (Password: employee123)
-- Set manager_id to the manager we just created
INSERT INTO users (email, password_hash, first_name, last_name, role, department, manager_id, is_active) VALUES
    ('employee@abccompany.com', '$2a$10$j3brMRlXOQdRfBQ/9adJ4em6HzYVzchEgSNKuGkwdipzLOolS.Xse', 'Employee', 'User', 'employee', 'Engineering',
    (SELECT id FROM users WHERE email = 'manager@abccompany.com'), true);

-- ============================================
-- Allocate Leave Balances for Test Users
-- Current Year: 2026
-- ============================================

-- Get current year
DO $$
DECLARE
    current_year INTEGER := EXTRACT(YEAR FROM CURRENT_DATE);
    admin_user_id INTEGER;
    manager_user_id INTEGER;
    employee_user_id INTEGER;
    annual_leave_id INTEGER;
    sick_leave_id INTEGER;
    casual_leave_id INTEGER;
BEGIN
    -- Get user IDs
    SELECT id INTO admin_user_id FROM users WHERE email = 'admin@abccompany.com';
    SELECT id INTO manager_user_id FROM users WHERE email = 'manager@abccompany.com';
    SELECT id INTO employee_user_id FROM users WHERE email = 'employee@abccompany.com';

    -- Get leave type IDs
    SELECT id INTO annual_leave_id FROM leave_types WHERE name = 'Annual Leave';
    SELECT id INTO sick_leave_id FROM leave_types WHERE name = 'Sick Leave';
    SELECT id INTO casual_leave_id FROM leave_types WHERE name = 'Casual Leave';

    -- Allocate leave balances for Admin
    INSERT INTO leave_balances (user_id, leave_type_id, total_days, used_days, year) VALUES
        (admin_user_id, annual_leave_id, 20, 0, current_year),
        (admin_user_id, sick_leave_id, 10, 0, current_year),
        (admin_user_id, casual_leave_id, 7, 0, current_year)
    ON CONFLICT (user_id, leave_type_id, year) DO UPDATE
        SET total_days = EXCLUDED.total_days;

    -- Allocate leave balances for Manager
    INSERT INTO leave_balances (user_id, leave_type_id, total_days, used_days, year) VALUES
        (manager_user_id, annual_leave_id, 20, 2, current_year),
        (manager_user_id, sick_leave_id, 10, 1, current_year),
        (manager_user_id, casual_leave_id, 7, 0, current_year)
    ON CONFLICT (user_id, leave_type_id, year) DO UPDATE
        SET total_days = EXCLUDED.total_days, used_days = EXCLUDED.used_days;

    -- Allocate leave balances for Employee
    INSERT INTO leave_balances (user_id, leave_type_id, total_days, used_days, year) VALUES
        (employee_user_id, annual_leave_id, 20, 5, current_year),
        (employee_user_id, sick_leave_id, 10, 2, current_year),
        (employee_user_id, casual_leave_id, 7, 1, current_year)
    ON CONFLICT (user_id, leave_type_id, year) DO UPDATE
        SET total_days = EXCLUDED.total_days, used_days = EXCLUDED.used_days;

    RAISE NOTICE 'Leave balances allocated for year %', current_year;
END $$;

-- ============================================
-- Create Sample Leave Requests for Testing
-- ============================================

DO $$
DECLARE
    employee_user_id INTEGER;
    manager_user_id INTEGER;
    annual_leave_id INTEGER;
    sick_leave_id INTEGER;
BEGIN
    -- Get IDs
    SELECT id INTO employee_user_id FROM users WHERE email = 'employee@abccompany.com';
    SELECT id INTO manager_user_id FROM users WHERE email = 'manager@abccompany.com';
    SELECT id INTO annual_leave_id FROM leave_types WHERE name = 'Annual Leave';
    SELECT id INTO sick_leave_id FROM leave_types WHERE name = 'Sick Leave';

    -- Sample pending leave request
    INSERT INTO leave_requests (user_id, leave_type_id, start_date, end_date, total_days, reason, status) VALUES
        (employee_user_id, annual_leave_id, CURRENT_DATE + INTERVAL '7 days', CURRENT_DATE + INTERVAL '9 days', 3, 'Family vacation', 'pending');

    -- Sample approved leave request
    INSERT INTO leave_requests (user_id, leave_type_id, start_date, end_date, total_days, reason, status, reviewed_by, reviewed_at, review_notes) VALUES
        (employee_user_id, sick_leave_id, CURRENT_DATE - INTERVAL '5 days', CURRENT_DATE - INTERVAL '3 days', 2, 'Medical appointment', 'approved', manager_user_id, CURRENT_TIMESTAMP - INTERVAL '4 days', 'Approved after medical certificate verification');

    -- Sample manager leave request
    INSERT INTO leave_requests (user_id, leave_type_id, start_date, end_date, total_days, reason, status) VALUES
        (manager_user_id, annual_leave_id, CURRENT_DATE + INTERVAL '14 days', CURRENT_DATE + INTERVAL '16 days', 3, 'Personal work', 'pending');

    RAISE NOTICE 'Sample leave requests created successfully';
END $$;

-- ============================================
-- Display Test Accounts Summary
-- ============================================

DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Test Accounts Created Successfully!';
    RAISE NOTICE '========================================';
    RAISE NOTICE '';
    RAISE NOTICE 'Admin Account:';
    RAISE NOTICE '  Email: admin@abccompany.com';
    RAISE NOTICE '  Password: admin123';
    RAISE NOTICE '  Role: admin';
    RAISE NOTICE '';
    RAISE NOTICE 'Manager Account:';
    RAISE NOTICE '  Email: manager@abccompany.com';
    RAISE NOTICE '  Password: manager123';
    RAISE NOTICE '  Role: manager';
    RAISE NOTICE '';
    RAISE NOTICE 'Employee Account:';
    RAISE NOTICE '  Email: employee@abccompany.com';
    RAISE NOTICE '  Password: employee123';
    RAISE NOTICE '  Role: employee';
    RAISE NOTICE '';
    RAISE NOTICE '========================================';
    RAISE NOTICE 'All accounts have leave balances allocated';
    RAISE NOTICE 'Sample leave requests created for testing';
    RAISE NOTICE '========================================';
END $$;
