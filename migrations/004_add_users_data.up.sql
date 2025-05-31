-- ลบข้อมูลเดิมก่อนเพิ่มข้อมูลใหม่
DELETE FROM users;

-- Reset sequence ของ ID users
ALTER SEQUENCE users_id_seq RESTART WITH 1;

-- เพิ่มข้อมูล users ตามรูปภาพที่ส่งมา
INSERT INTO users (username, password, email, provider, role) VALUES 
('admin', 'password', 'admin@email.com', 'local', 'admin'),
('user1', 'password', 'user1@example.com', 'local', 'user'),
('user2', 'password', 'user2@example.com', 'local', 'user');
