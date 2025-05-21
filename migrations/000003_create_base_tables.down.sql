-- Drop tables in reverse order (to handle foreign key constraints)
DROP TABLE IF EXISTS favorites;
DROP TABLE IF EXISTS oauth_tokens;
DROP TABLE IF EXISTS schedule_logs;
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS staff;
DROP TABLE IF EXISTS vehicles;
DROP TABLE IF EXISTS routes;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS stations;