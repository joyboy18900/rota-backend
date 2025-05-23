-- ลบคอลัมน์ deleted_at
ALTER TABLE schedules 
  DROP COLUMN IF EXISTS deleted_at;

-- ลบคอลัมน์ status
ALTER TABLE schedules 
  DROP COLUMN IF EXISTS status;

-- ลบคอลัมน์ vehicle_id
ALTER TABLE schedules 
  DROP COLUMN IF EXISTS vehicle_id;

-- เปลี่ยนประเภทข้อมูลกลับเป็น time without time zone
ALTER TABLE schedules 
  ALTER COLUMN departure_time TYPE time without time zone,
  ALTER COLUMN arrival_time TYPE time without time zone;
