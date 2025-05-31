-- ลบข้อมูลที่เพิ่มในไฟล์ up migration
DELETE FROM schedule_logs;
DELETE FROM schedules;
DELETE FROM vehicles;
DELETE FROM routes;
DELETE FROM stations;
