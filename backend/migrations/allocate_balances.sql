-- Allocate Leave Balances for Test Users
-- Year: 2026

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
    ON CONFLICT (user_id, leave_type_id, year) DO UPDATE
        SET total_days = EXCLUDED.total_days;

    RAISE NOTICE 'Leave balances allocated successfully for year %', current_year;
END $$;
