-- Allocate Leave Balances for Test Users
-- Year: 2026

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
        SET total_days = EXCLUDED.total_days, used_days = EXCLUDED.used_days;

    -- Allocate leave balances for Manager
    INSERT INTO leave_balances (user_id, leave_type_id, total_days, used_days, year) VALUES
        (manager_user_id, annual_leave_id, 20, 0, current_year),
        (manager_user_id, sick_leave_id, 10, 0, current_year),
        (manager_user_id, casual_leave_id, 7, 0, current_year)
    ON CONFLICT (user_id, leave_type_id, year) DO UPDATE
        SET total_days = EXCLUDED.total_days, used_days = EXCLUDED.used_days;

    -- Allocate leave balances for Employee
    INSERT INTO leave_balances (user_id, leave_type_id, total_days, used_days, year) VALUES
        (employee_user_id, annual_leave_id, 20, 0, current_year),
        (employee_user_id, sick_leave_id, 10, 0, current_year),
        (employee_user_id, casual_leave_id, 7, 0, current_year)
    ON CONFLICT (user_id, leave_type_id, year) DO UPDATE
        SET total_days = EXCLUDED.total_days, used_days = EXCLUDED.used_days;

    RAISE NOTICE 'Leave balances allocated successfully for year %', current_year;
END $$;
