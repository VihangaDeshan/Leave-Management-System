-- Seed leave types
INSERT INTO leave_types (name, description, is_active) VALUES
('Annual Leave', 'Yearly vacation leave', true),
('Sick Leave', 'Medical or health-related leave', true),
('Casual Leave', 'Short-term personal leave', true),
('Maternity Leave', 'Leave for mothers before and after childbirth', true),
('Paternity Leave', 'Leave for fathers after childbirth', true),
('Unpaid Leave', 'Leave without pay', true),
('Bereavement Leave', 'Leave due to death of a family member', true)
ON CONFLICT (name) DO NOTHING;

-- Create default admin user
-- Password: admin123 (hashed with bcrypt, cost 12)
-- Note: Change this password after first login in production!
INSERT INTO users (email, password_hash, first_name, last_name, role, department, is_active) VALUES
('admin@abccompany.com', '$2a$12$LQv3c1yqBW.PFcHdNKr0AeFYKwvfbCJvF9ZL4lNLdXsG2KxkVL.YK', 'System', 'Administrator', 'admin', 'Administration', true)
ON CONFLICT (email) DO NOTHING;

-- Create sample manager user
-- Password: manager123 (hashed with bcrypt, cost 12)
INSERT INTO users (email, password_hash, first_name, last_name, role, department, is_active) VALUES
('manager@abccompany.com', '$2a$12$4.RyfLKKBvq0j9l5pKxz5uKJB4X7r7wBZL8vXwX.5xP8qS.K9nK7C', 'John', 'Manager', 'manager', 'Human Resources', true)
ON CONFLICT (email) DO NOTHING;

-- Create sample employee user
-- Password: employee123 (hashed with bcrypt, cost 12)
INSERT INTO users (email, password_hash, first_name, last_name, role, department, is_active) VALUES
('employee@abccompany.com', '$2a$12$8.TyfLKKBvq0j9l5pKxz5uKJB4X7r7wBZL8vXwX.5xP8qS.K9nK8D', 'Jane', 'Employee', 'employee', 'Engineering', true)
ON CONFLICT (email) DO NOTHING;

-- Allocate leave balances for current year for sample users
-- Get current year dynamically
DO $$
DECLARE
    current_year INTEGER := EXTRACT(YEAR FROM CURRENT_DATE);
    admin_user_id INTEGER;
    manager_user_id INTEGER;
    employee_user_id INTEGER;
    annual_leave_type_id INTEGER;
    sick_leave_type_id INTEGER;
    casual_leave_type_id INTEGER;
BEGIN
    -- Get user IDs
    SELECT id INTO admin_user_id FROM users WHERE email = 'admin@abccompany.com';
    SELECT id INTO manager_user_id FROM users WHERE email = 'manager@abccompany.com';
    SELECT id INTO employee_user_id FROM users WHERE email = 'employee@abccompany.com';

    -- Get leave type IDs
    SELECT id INTO annual_leave_type_id FROM leave_types WHERE name = 'Annual Leave';
    SELECT id INTO sick_leave_type_id FROM leave_types WHERE name = 'Sick Leave';
    SELECT id INTO casual_leave_type_id FROM leave_types WHERE name = 'Casual Leave';

    -- Allocate leave balances for admin
    IF admin_user_id IS NOT NULL THEN
        INSERT INTO leave_balances (user_id, leave_type_id, total_days, used_days, year)
        VALUES
            (admin_user_id, annual_leave_type_id, 20, 0, current_year),
            (admin_user_id, sick_leave_type_id, 10, 0, current_year),
            (admin_user_id, casual_leave_type_id, 5, 0, current_year)
        ON CONFLICT (user_id, leave_type_id, year) DO NOTHING;
    END IF;

    -- Allocate leave balances for manager
    IF manager_user_id IS NOT NULL THEN
        INSERT INTO leave_balances (user_id, leave_type_id, total_days, used_days, year)
        VALUES
            (manager_user_id, annual_leave_type_id, 20, 0, current_year),
            (manager_user_id, sick_leave_type_id, 10, 0, current_year),
            (manager_user_id, casual_leave_type_id, 5, 0, current_year)
        ON CONFLICT (user_id, leave_type_id, year) DO NOTHING;
    END IF;

    -- Allocate leave balances for employee
    IF employee_user_id IS NOT NULL THEN
        INSERT INTO leave_balances (user_id, leave_type_id, total_days, used_days, year)
        VALUES
            (employee_user_id, annual_leave_type_id, 20, 0, current_year),
            (employee_user_id, sick_leave_type_id, 10, 0, current_year),
            (employee_user_id, casual_leave_type_id, 5, 0, current_year)
        ON CONFLICT (user_id, leave_type_id, year) DO NOTHING;
    END IF;
END $$;
