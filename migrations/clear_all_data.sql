-- ลบข้อมูลทั้งหมดตามลำดับ foreign key constraints
-- ต้องลบข้อมูลจากตาราง favorites ก่อนเพราะมี foreign key ไปยังตารางอื่น
DELETE FROM favorites;
DELETE FROM schedule_logs;
DELETE FROM schedules;
DELETE FROM vehicles;
DELETE FROM routes;
DELETE FROM stations;
DELETE FROM users CASCADE;
DELETE FROM staffs CASCADE;

-- Reset sequence ของ ID ทุกตาราง
ALTER SEQUENCE stations_id_seq RESTART WITH 1;
ALTER SEQUENCE routes_id_seq RESTART WITH 1;
ALTER SEQUENCE vehicles_id_seq RESTART WITH 1;
ALTER SEQUENCE schedules_id_seq RESTART WITH 1;
ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE staffs_id_seq RESTART WITH 1;
