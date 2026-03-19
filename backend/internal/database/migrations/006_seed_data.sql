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
BEGIN
    INSERT INTO leave_balances (user_id, leave_type_id, total_days, used_days, year)
    SELECT
        u.id,
        lt.id,
        CASE
            WHEN lower(lt.name) = 'annual leave' THEN 20
            WHEN lower(lt.name) = 'sick leave' THEN 10
            WHEN lower(lt.name) = 'casual leave' THEN 5
            ELSE 0
        END AS total_days,
        0,
        current_year
    FROM users u
    CROSS JOIN leave_types lt
    WHERE u.email IN ('admin@abccompany.com', 'manager@abccompany.com', 'employee@abccompany.com')
      AND lt.is_active = true
    ON CONFLICT (user_id, leave_type_id, year) DO NOTHING;
END $$;
