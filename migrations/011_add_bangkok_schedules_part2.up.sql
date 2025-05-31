-- สร้างตารางเวลาเดินรถสำหรับเส้นทาง หมอชิต 2 <-> มีนบุรี
-- 1. หมอชิต 2 -> มีนบุรี
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status) 
SELECT 
  7 AS route_id, 
  7 AS vehicle_id,
  1 AS station_id,
  generate_series AS round,
  ('2025-06-01 06:45:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 07:20:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);

-- 2. มีนบุรี -> หมอชิต 2
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status)
SELECT 
  8 AS route_id, 
  8 AS vehicle_id,
  5 AS station_id,
  generate_series AS round,
  ('2025-06-01 07:45:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 08:20:00'::timestamp + ((generate_series - 1) * interval '1 hour 30 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);

-- สร้างตารางเวลาเดินรถสำหรับเส้นทาง หมอชิต 2 <-> หมอชิตเก่า
-- 1. หมอชิต 2 -> หมอชิตเก่า
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status) 
SELECT 
  9 AS route_id, 
  9 AS vehicle_id,
  1 AS station_id,
  generate_series AS round,
  ('2025-06-01 06:00:00'::timestamp + ((generate_series - 1) * interval '1 hour'))::timestamptz AS departure_time,
  ('2025-06-01 06:10:00'::timestamp + ((generate_series - 1) * interval '1 hour'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);

-- 2. หมอชิตเก่า -> หมอชิต 2
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status)
SELECT 
  10 AS route_id, 
  10 AS vehicle_id,
  6 AS station_id,
  generate_series AS round,
  ('2025-06-01 06:30:00'::timestamp + ((generate_series - 1) * interval '1 hour'))::timestamptz AS departure_time,
  ('2025-06-01 06:40:00'::timestamp + ((generate_series - 1) * interval '1 hour'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);

-- สร้างตารางเวลาเดินรถสำหรับเส้นทาง หมอชิต 2 <-> บางเขน
-- 1. หมอชิต 2 -> บางเขน
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status) 
SELECT 
  11 AS route_id, 
  11 AS vehicle_id,
  1 AS station_id,
  generate_series AS round,
  ('2025-06-01 06:20:00'::timestamp + ((generate_series - 1) * interval '1 hour 15 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 06:35:00'::timestamp + ((generate_series - 1) * interval '1 hour 15 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);

-- 2. บางเขน -> หมอชิต 2
INSERT INTO schedules (route_id, vehicle_id, station_id, round, departure_time, arrival_time, status)
SELECT 
  12 AS route_id, 
  12 AS vehicle_id,
  7 AS station_id,
  generate_series AS round,
  ('2025-06-01 06:50:00'::timestamp + ((generate_series - 1) * interval '1 hour 15 minutes'))::timestamptz AS departure_time,
  ('2025-06-01 07:05:00'::timestamp + ((generate_series - 1) * interval '1 hour 15 minutes'))::timestamptz AS arrival_time,
  'scheduled' AS status
FROM generate_series(1, 10);
