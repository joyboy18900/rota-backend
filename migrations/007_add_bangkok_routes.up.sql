-- สร้างเส้นทางระหว่างหมอชิต 2 กับสถานีอื่นๆ ในกรุงเทพฯ
INSERT INTO routes (start_station_id, end_station_id, distance, duration) VALUES
-- หมอชิต 2 <-> รังสิต
(1, 2, 22.5, '30m'),
(2, 1, 22.5, '30m'),
-- หมอชิต 2 <-> สายใต้ใหม่
(1, 3, 26.8, '40m'),
(3, 1, 26.8, '40m'),
-- หมอชิต 2 <-> บางนา
(1, 4, 25.3, '35m'),
(4, 1, 25.3, '35m'),
-- หมอชิต 2 <-> มีนบุรี
(1, 5, 22.0, '35m'),
(5, 1, 22.0, '35m'),
-- หมอชิต 2 <-> หมอชิตเก่า
(1, 6, 2.5, '10m'),
(6, 1, 2.5, '10m'),
-- หมอชิต 2 <-> บางเขน
(1, 7, 8.7, '15m'),
(7, 1, 8.7, '15m'),
-- หมอชิต 2 <-> ลาดพร้าว
(1, 8, 9.2, '20m'),
(8, 1, 9.2, '20m'),
-- หมอชิต 2 <-> แจ้งวัฒนะ
(1, 9, 16.1, '25m'),
(9, 1, 16.1, '25m'),
-- หมอชิต 2 <-> อนุสาวรีย์ชัย
(1, 10, 7.3, '15m'),
(10, 1, 7.3, '15m');
