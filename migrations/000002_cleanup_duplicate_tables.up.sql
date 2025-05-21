-- First, drop foreign key constraints
ALTER TABLE IF EXISTS schedule_logs
    DROP CONSTRAINT IF EXISTS schedule_logs_staff_id_fkey;

-- Drop duplicate tables if they exist
DROP TABLE IF EXISTS staff CASCADE;
DROP TABLE IF EXISTS o_auth_tokens CASCADE;

-- Rename staffs to staff if needed
ALTER TABLE IF EXISTS staffs RENAME TO staff;

-- Rename oauth_tokens to o_auth_tokens if needed
ALTER TABLE IF EXISTS oauth_tokens RENAME TO o_auth_tokens;

-- Recreate foreign key constraints
ALTER TABLE IF EXISTS schedule_logs
    ADD CONSTRAINT schedule_logs_staff_id_fkey
    FOREIGN KEY (staff_id)
    REFERENCES staff(id);