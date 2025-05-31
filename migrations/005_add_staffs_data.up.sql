-- ลบข้อมูล staffs เดิมก่อนเพิ่มข้อมูลใหม่
DELETE FROM staffs;

-- Reset sequence ของ ID staffs
ALTER SEQUENCE staffs_id_seq RESTART WITH 1;

-- เพิ่มข้อมูล staffs ตามรูปภาพที่ส่งมา
INSERT INTO staffs (username, password, email, name, position) VALUES 
('somchai', 'password', 'somchai@example.com', 'Somchai Jaidee', 'admin'),
('somsak', 'password', 'somsak@example.com', 'Somsak Jaidee', 'admin'),
('staff1', 'password', 'staff1@example.com', 'พนักงาน ทดสอบ (อัพเดท)', 'staff'),
('staff2', 'password', 'staff2@example.com', 'Test Staff', 'staff'),
('staff3', 'password', 'staff3@example.com', 'Staff Test', 'staff');
