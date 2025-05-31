-- ลบข้อมูลตารางเวลาเดินรถที่เพิ่มใน up migration
DELETE FROM schedules WHERE route_id IN (13, 14, 15, 16, 17, 18);
