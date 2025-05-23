-- เปลี่ยนประเภทข้อมูลของ departure_time และ arrival_time
ALTER TABLE schedules 
  ALTER COLUMN departure_time TYPE timestamp with time zone,
  ALTER COLUMN arrival_time TYPE timestamp with time zone;

-- เพิ่มคอลัมน์ vehicle_id
ALTER TABLE schedules 
  ADD COLUMN vehicle_id INTEGER REFERENCES vehicles(id);

-- เพิ่มคอลัมน์ status
ALTER TABLE schedules 
  ADD COLUMN status VARCHAR(20) DEFAULT 'scheduled';

-- เพิ่มคอลัมน์ deleted_at สำหรับ soft delete
ALTER TABLE schedules 
  ADD COLUMN deleted_at timestamp with time zone;
