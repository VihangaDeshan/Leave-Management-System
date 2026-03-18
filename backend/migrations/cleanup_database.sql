-- Cleanup Script: Reset database to clean state
-- Keeps: leave_types and test accounts only
-- Removes: all leave_requests and leave_balances

-- ============================================
-- Truncate Tables (Clear Data)
-- ============================================

-- Delete all leave requests
TRUNCATE TABLE leave_requests CASCADE;

-- Delete all leave balances
TRUNCATE TABLE leave_balances CASCADE;

-- Delete all audit logs
TRUNCATE TABLE audit_logs CASCADE;

-- Keep leave_types and users (they will be preserved)

-- ============================================
-- Success Message
-- ============================================
DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Database Cleaned Successfully!';
    RAISE NOTICE '========================================';
    RAISE NOTICE '';
    RAISE NOTICE 'Preserved:';
    RAISE NOTICE '  - Leave Types (6 types)';
    RAISE NOTICE '  - Test User Accounts (3 users)';
    RAISE NOTICE '';
    RAISE NOTICE 'Removed:';
    RAISE NOTICE '  - All Leave Requests';
    RAISE NOTICE '  - All Leave Balances';
    RAISE NOTICE '  - All Audit Logs';
    RAISE NOTICE '';
    RAISE NOTICE '========================================';
END $$;
