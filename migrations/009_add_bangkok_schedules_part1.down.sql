-- ลบข้อมูลตารางเวลาเดินรถที่เพิ่มใน up migration
DELETE FROM schedules WHERE route_id IN (1, 2, 3, 4, 5, 6);
