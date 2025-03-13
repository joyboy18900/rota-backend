-- Users (ตัวอย่างผู้ใช้งาน)
INSERT INTO users (username, password, email, provider, is_verified)
VALUES
('user1', 'password1', 'user1@example.com', 'local', 1),
('user2', NULL, 'user2@example.com', 'google', 1),
('user3', 'password3', 'user3@example.com', 'local', 0),
('user4', NULL, 'user4@example.com', 'google', 1),
('user5', 'password5', 'user5@example.com', 'local', 1);

-- Stations (ตัวอย่างสถานี)
INSERT INTO stations (name, location)
VALUES
('หมอชิต', 'กรุงเทพฯ'),
('เชียงใหม่', 'เชียงใหม่'),
('นครราชสีมา', 'โคราช'),
('ขอนแก่น', 'ขอนแก่น'),
('อุบลราชธานี', 'อุบลฯ'),
('หาดใหญ่', 'สงขลา'),
('ภูเก็ต', 'ภูเก็ต'),
('พัทยา', 'ชลบุรี'),
('ระยอง', 'ระยอง'),
('เชียงราย', 'เชียงราย');

-- Routes (ตัวอย่างเส้นทาง)
INSERT INTO routes (start_station_id, end_station_id, distance, duration)
VALUES
(1, 2, 700, '10:00'),
(1, 3, 260, '04:00'),
(1, 4, 450, '06:00'),
(1, 5, 630, '08:00'),
(1, 6, 950, '12:00'),
(1, 7, 860, '11:00'),
(1, 8, 150, '02:30'),
(1, 9, 220, '03:30'),
(1, 10, 800, '10:00'),
(2, 5, 300, '04:30');

-- Schedules (ตารางเดินรถ)
INSERT INTO schedules (route_id, station_id, round, departure_time, arrival_time)
VALUES
(1, 1, 1, '06:00', '16:00'),
(1, 2, 2, '08:00', '18:00'),
(2, 1, 1, '07:00', '11:00'),
(2, 3, 2, '09:00', '13:00'),
(3, 1, 1, '05:30', '11:30'),
(3, 4, 2, '07:30', '13:30'),
(4, 1, 1, '06:00', '14:00'),
(4, 5, 2, '09:00', '17:00'),
(5, 1, 1, '04:00', '16:00'),
(5, 6, 2, '08:00', '20:00');

-- Staff (เจ้าหน้าที่)
INSERT INTO staff (username, password, email, station_id)
VALUES
('staff1', 'staffpass1', 'staff1@example.com', 1),
('staff2', 'staffpass2', 'staff2@example.com', 2),
('staff3', 'staffpass3', 'staff3@example.com', 3),
('staff4', 'staffpass4', 'staff4@example.com', 4),
('staff5', 'staffpass5', 'staff5@example.com', 5);

-- Vehicles (ข้อมูลรถตู้)
INSERT INTO vehicles (license_plate, capacity, driver_name, route_id)
VALUES
('บข 1234', 15, 'สมชาย ขับดี', 1),
('บข 5678', 12, 'สมศรี ใจดี', 2),
('บข 9101', 18, 'สมปอง พาทัวร์', 3),
('บข 1121', 15, 'สมศักดิ์ ขับไว', 4),
('บข 3141', 12, 'สมหมาย เร็วแรง', 5);

-- OAuth Tokens (ข้อมูล Token สำหรับ OAuth)
INSERT INTO oauth_tokens (user_id, provider, access_token, refresh_token, expires_at)
VALUES
(3, 'google', 'token_google_1', 'refresh_google_1', '2024-12-31 23:59:59'),
(5, 'google', 'token_google_2', 'refresh_google_2', '2024-12-31 23:59:59');

-- Favorites (ข้อมูลสถานีที่ผู้ใช้งานชื่นชอบ)
INSERT INTO favorites (user_id, station_id, created_at)
VALUES
(1, 1, CURRENT_TIMESTAMP),
(1, 3, CURRENT_TIMESTAMP),
(2, 5, CURRENT_TIMESTAMP),
(3, 2, CURRENT_TIMESTAMP),
(4, 4, CURRENT_TIMESTAMP),
(5, 1, CURRENT_TIMESTAMP),
(1, 6, CURRENT_TIMESTAMP),
(2, 8, CURRENT_TIMESTAMP),
(3, 10, CURRENT_TIMESTAMP),
(4, 9, CURRENT_TIMESTAMP);

-- Schedule Logs (บันทึกการเปลี่ยนแปลงตารางเดินรถ)
INSERT INTO schedule_logs (schedule_id, staff_id, change_description, updated_at)
VALUES
(1, 1, 'เปลี่ยนเวลาออกเดินทางจาก 06:00 เป็น 06:30', CURRENT_TIMESTAMP),
(2, 2, 'เพิ่มรอบเดินรถรอบใหม่ 09:00', CURRENT_TIMESTAMP),
(3, 3, 'แก้ไขเวลาถึงปลายทางจาก 13:00 เป็น 12:45', CURRENT_TIMESTAMP),
(4, 4, 'ลบตารางเดินรถที่ไม่ได้ใช้งาน', CURRENT_TIMESTAMP),
(5, 5, 'เปลี่ยนรอบเดินรถจาก 1 เป็น 2', CURRENT_TIMESTAMP),
(6, 1, 'แก้ไขเวลารถออกเป็น 07:15', CURRENT_TIMESTAMP),
(7, 2, 'เปลี่ยนเส้นทางเดินรถไปยังสถานีอื่น', CURRENT_TIMESTAMP),
(8, 3, 'เพิ่มรอบพิเศษช่วงวันหยุด', CURRENT_TIMESTAMP),
(9, 4, 'เปลี่ยนเวลาถึงปลายทางจาก 17:00 เป็น 16:50', CURRENT_TIMESTAMP),
(10, 5, 'ลดจำนวนรอบวิ่งลง 1 รอบ', CURRENT_TIMESTAMP);

-- Vehicle Logs (บันทึกการใช้งานรถตู้)
INSERT INTO vehicles (license_plate, capacity, driver_name, route_id)
VALUES
('บข 5555', 15, 'สมปอง มั่นใจ', 6),
('บข 6666', 18, 'สมพร ซิ่งแรง', 7),
('บข 7777', 12, 'สมชาย สุขใจ', 8),
('บข 8888', 15, 'สมศรี คนดี', 9),
('บข 9999', 20, 'สมหมาย พาเพลิน', 10);
