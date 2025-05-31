-- ลบข้อมูลตารางเวลาเดินรถที่เพิ่มใน up migration
DELETE FROM schedules WHERE route_id IN (7, 8, 9, 10, 11, 12);
