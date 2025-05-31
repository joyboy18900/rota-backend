-- สร้างตารางเวลาเดินรถสำหรับเส้นทาง หมอชิต 2 <-> ลาดพร้าว
-- 1. หมอชิต 2 -> ลาดพร้าว
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status) 
SELECT 
  13 AS route_id, 
  13 AS vehicle_id,
  1 AS station_id,
  generate_series AS round,
  ('2025-06-01 06:30:00'::timestamp + ((generate_series - 1) * interval '1 hour 20 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 06:50:00'::timestamp + ((generate_series - 1) * interval '1 hour 20 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);

-- 2. ลาดพร้าว -> หมอชิต 2
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status)
SELECT 
  14 AS route_id, 
  14 AS vehicle_id,
  8 AS station_id,
  generate_series AS round,
  ('2025-06-01 07:00:00'::timestamp + ((generate_series - 1) * interval '1 hour 20 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 07:20:00'::timestamp + ((generate_series - 1) * interval '1 hour 20 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);

-- สร้างตารางเวลาเดินรถสำหรับเส้นทาง หมอชิต 2 <-> แจ้งวัฒนะ
-- 1. หมอชิต 2 -> แจ้งวัฒนะ
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status) 
SELECT 
  15 AS route_id, 
  15 AS vehicle_id,
  1 AS station_id,
  generate_series AS round,
  ('2025-06-01 06:15:00'::timestamp + ((generate_series - 1) * interval '1 hour 25 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 06:40:00'::timestamp + ((generate_series - 1) * interval '1 hour 25 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);

-- 2. แจ้งวัฒนะ -> หมอชิต 2
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status)
SELECT 
  16 AS route_id, 
  16 AS vehicle_id,
  9 AS station_id,
  generate_series AS round,
  ('2025-06-01 07:15:00'::timestamp + ((generate_series - 1) * interval '1 hour 25 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 07:40:00'::timestamp + ((generate_series - 1) * interval '1 hour 25 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);

-- สร้างตารางเวลาเดินรถสำหรับเส้นทาง หมอชิต 2 <-> อนุสาวรีย์ชัย
-- 1. หมอชิต 2 -> อนุสาวรีย์ชัย
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status) 
SELECT 
  17 AS route_id, 
  17 AS vehicle_id,
  1 AS station_id,
  generate_series AS round,
  ('2025-06-01 06:00:00'::timestamp + ((generate_series - 1) * interval '1 hour 15 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 06:15:00'::timestamp + ((generate_series - 1) * interval '1 hour 15 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);

-- 2. อนุสาวรีย์ชัย -> หมอชิต 2
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status)
SELECT 
  18 AS route_id, 
  18 AS vehicle_id,
  10 AS station_id,
  generate_series AS round,
  ('2025-06-01 06:30:00'::timestamp + ((generate_series - 1) * interval '1 hour 15 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 06:45:00'::timestamp + ((generate_series - 1) * interval '1 hour 15 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);
