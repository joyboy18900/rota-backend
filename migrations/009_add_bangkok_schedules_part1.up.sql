-- สร้างตารางเวลาเดินรถสำหรับเส้นทาง หมอชิต 2 <-> รังสิต
-- 1. หมอชิต 2 -> รังสิต (10 รอบต่อวัน เริ่มตั้งแต่ 6:00 - 20:00 ทุกๆ 1.5 ชั่วโมง)
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status) 
SELECT 
  1 AS route_id, 
  1 AS vehicle_id,
  1 AS station_id,
  generate_series AS round,
  ('2025-06-01 06:00:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 06:30:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);

-- 2. รังสิต -> หมอชิต 2 (10 รอบต่อวัน เริ่มตั้งแต่ 7:00 - 21:00 ทุกๆ 1.5 ชั่วโมง)
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status)
SELECT 
  2 AS route_id, 
  2 AS vehicle_id,
  2 AS station_id,
  generate_series AS round,
  ('2025-06-01 07:00:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 07:30:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);

-- สร้างตารางเวลาเดินรถสำหรับเส้นทาง หมอชิต 2 <-> สายใต้ใหม่
-- 1. หมอชิต 2 -> สายใต้ใหม่ (10 รอบต่อวัน เริ่มตั้งแต่ 5:30 - 19:30 ทุกๆ 1.5 ชั่วโมง)
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status) 
SELECT 
  3 AS route_id, 
  3 AS vehicle_id,
  1 AS station_id,
  generate_series AS round,
  ('2025-06-01 05:30:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 06:10:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);

-- 2. สายใต้ใหม่ -> หมอชิต 2 (10 รอบต่อวัน เริ่มตั้งแต่ 6:30 - 20:30 ทุกๆ 1.5 ชั่วโมง)
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status)
SELECT 
  4 AS route_id, 
  4 AS vehicle_id,
  3 AS station_id,
  generate_series AS round,
  ('2025-06-01 06:30:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 07:10:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);

-- สร้างตารางเวลาเดินรถสำหรับเส้นทาง หมอชิต 2 <-> บางนา
-- 1. หมอชิต 2 -> บางนา (10 รอบต่อวัน เริ่มตั้งแต่ 6:15 - 20:15 ทุกๆ 1.5 ชั่วโมง)
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status) 
SELECT 
  5 AS route_id, 
  5 AS vehicle_id,
  1 AS station_id,
  generate_series AS round,
  ('2025-06-01 06:15:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 06:50:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);

-- 2. บางนา -> หมอชิต 2 (10 รอบต่อวัน เริ่มตั้งแต่ 7:15 - 21:15 ทุกๆ 1.5 ชั่วโมง)
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status)
SELECT 
  6 AS route_id, 
  6 AS vehicle_id,
  4 AS station_id,
  generate_series AS round,
  ('2025-06-01 07:15:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 07:50:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);
