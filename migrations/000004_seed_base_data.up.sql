-- Seed stations
INSERT INTO stations (id, name, description, latitude, longitude, address, status) VALUES
    ('11111111-1111-1111-1111-111111111111', 'มหาวิทยาลัยเชียงใหม่', 'สถานีหน้ามหาวิทยาลัยเชียงใหม่', 18.796143, 98.952617, 'ถ.ห้วยแก้ว ต.สุเทพ อ.เมือง จ.เชียงใหม่', 'active'),
    ('22222222-2222-2222-2222-222222222222', 'ตลาดวโรรส', 'สถานีกาดหลวง', 18.787677, 98.997836, 'ถ.วิชยานนท์ ต.ช้างม่อย อ.เมือง จ.เชียงใหม่', 'active'),
    ('33333333-3333-3333-3333-333333333333', 'เซ็นทรัลเชียงใหม่แอร์พอร์ต', 'สถานีหน้าห้างเซ็นทรัล', 18.766869, 98.974595, 'ถ.มหิดล ต.หายยา อ.เมือง จ.เชียงใหม่', 'active'),
    ('44444444-4444-4444-4444-444444444444', 'นิมมานเหมินทร์', 'สถานีถนนนิมมานเหมินทร์', 18.798787, 98.967400, 'ถ.นิมมานเหมินทร์ ต.สุเทพ อ.เมือง จ.เชียงใหม่', 'active'),
    ('55555555-5555-5555-5555-555555555555', 'ท่าแพ', 'สถานีประตูท่าแพ', 18.787570, 98.993109, 'ถ.ท่าแพ ต.ช้างคลาน อ.เมือง จ.เชียงใหม่', 'active'),
    ('66666666-6666-6666-6666-666666666666', 'สนามบินเชียงใหม่', 'สถานีสนามบินนานาชาติเชียงใหม่', 18.766869, 98.962707, 'ถ.มหิดล ต.สุเทพ อ.เมือง จ.เชียงใหม่', 'active'),
    ('77777777-7777-7777-7777-777777777777', 'มีโชคพลาซ่า', 'สถานีห้างมีโชคพลาซ่า', 18.802458, 99.017431, 'ถ.เชียงใหม่-ลำปาง ต.ฟ้าฮ่าม อ.เมือง จ.เชียงใหม่', 'active'),
    ('88888888-8888-8888-8888-888888888888', 'เมญ่า', 'สถานีห้างเมญ่า', 18.802816, 98.967267, 'ถ.ห้วยแก้ว ต.สุเทพ อ.เมือง จ.เชียงใหม่', 'active'),
    ('99999999-9999-9999-9999-999999999999', 'กาดสวนแก้ว', 'สถานีห้างกาดสวนแก้ว', 18.797697, 98.977636, 'ถ.ห้วยแก้ว ต.สุเทพ อ.เมือง จ.เชียงใหม่', 'active'),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'ตลาดต้นพยอม', 'สถานีกาดต้นพยอม', 18.795459, 98.960464, 'ถ.สุเทพ ต.สุเทพ อ.เมือง จ.เชียงใหม่', 'active');

-- Seed routes
INSERT INTO routes (id, name, description, status) VALUES
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'สาย A: มช. - ตลาดวโรรส', 'เส้นทางจากมหาวิทยาลัยเชียงใหม่ไปตลาดวโรรส', 'active'),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', 'สาย B: เซ็นทรัล - สนามบิน', 'เส้นทางจากเซ็นทรัลไปสนามบิน', 'active'),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', 'สาย C: นิมมาน - มีโชค', 'เส้นทางจากนิมมานไปมีโชคพลาซ่า', 'active'),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'สาย D: ท่าแพ - เมญ่า', 'เส้นทางจากประตูท่าแพไปห้างเมญ่า', 'active'),
    ('ffffffff-ffff-ffff-ffff-ffffffffffff', 'สาย E: วนรอบเมือง', 'เส้นทางวนรอบเมืองเชียงใหม่', 'active');

-- Seed vehicles
INSERT INTO vehicles (id, name, vehicle_number, type, capacity, status) VALUES
    ('11111111-2222-3333-4444-555555555555', 'รถบัส A1', 'CMU-001', 'bus', 40, 'active'),
    ('22222222-3333-4444-5555-666666666666', 'รถบัส A2', 'CMU-002', 'bus', 40, 'active'),
    ('33333333-4444-5555-6666-777777777777', 'รถบัส B1', 'CMU-003', 'bus', 40, 'active'),
    ('44444444-5555-6666-7777-888888888888', 'รถบัส B2', 'CMU-004', 'bus', 40, 'active'),
    ('55555555-6666-7777-8888-999999999999', 'รถตู้ C1', 'CMU-005', 'van', 14, 'active'),
    ('66666666-7777-8888-9999-aaaaaaaaaaaa', 'รถตู้ C2', 'CMU-006', 'van', 14, 'active'),
    ('77777777-8888-9999-aaaa-bbbbbbbbbbbb', 'รถตู้ D1', 'CMU-007', 'van', 14, 'active'),
    ('88888888-9999-aaaa-bbbb-cccccccccccc', 'รถตู้ D2', 'CMU-008', 'van', 14, 'active'),
    ('99999999-aaaa-bbbb-cccc-dddddddddddd', 'รถบัส E1', 'CMU-009', 'bus', 40, 'active'),
    ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', 'รถบัส E2', 'CMU-010', 'bus', 40, 'active');

-- Seed route_stations
INSERT INTO route_stations (route_id, station_id, sequence_number) VALUES
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '11111111-1111-1111-1111-111111111111', 1),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '22222222-2222-2222-2222-222222222222', 2),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', '33333333-3333-3333-3333-333333333333', 1),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', '66666666-6666-6666-6666-666666666666', 2),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', '44444444-4444-4444-4444-444444444444', 1),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', '77777777-7777-7777-7777-777777777777', 2),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '55555555-5555-5555-5555-555555555555', 1),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '88888888-8888-8888-8888-888888888888', 2),
    ('ffffffff-ffff-ffff-ffff-ffffffffffff', '99999999-9999-9999-9999-999999999999', 1),
    ('ffffffff-ffff-ffff-ffff-ffffffffffff', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 2);

-- Seed route_vehicles
INSERT INTO route_vehicles (route_id, vehicle_id) VALUES
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '11111111-2222-3333-4444-555555555555'),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '22222222-3333-4444-5555-666666666666'),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', '33333333-4444-5555-6666-777777777777'),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', '44444444-5555-6666-7777-888888888888'),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', '55555555-6666-7777-8888-999999999999'),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', '66666666-7777-8888-9999-aaaaaaaaaaaa'),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '77777777-8888-9999-aaaa-bbbbbbbbbbbb'),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '88888888-9999-aaaa-bbbb-cccccccccccc'),
    ('ffffffff-ffff-ffff-ffff-ffffffffffff', '99999999-aaaa-bbbb-cccc-dddddddddddd'),
    ('ffffffff-ffff-ffff-ffff-ffffffffffff', 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');

-- Insert admin user
INSERT INTO users (id, username, email, password, role, is_active) VALUES
    ('admin1111-1111-1111-1111-111111111111', 'admin', 'admin@rota.local', '$2a$10$EqGca1YJP8GHmXKL1Z.kEuZa6KJL.Xj8R1XhYYqgwqXhvZzEXpYf.', 'admin', true);

-- Insert staff members
INSERT INTO staff (id, username, password, email, station_id) VALUES
    ('staff1111-1111-1111-1111-111111111111', 'staff_mcu', '$2a$10$EqGca1YJP8GHmXKL1Z.kEuZa6KJL.Xj8R1XhYYqgwqXhvZzEXpYf.', 'staff_mcu@rota.local', '11111111-1111-1111-1111-111111111111'),
    ('staff2222-2222-2222-2222-222222222222', 'staff_warorot', '$2a$10$EqGca1YJP8GHmXKL1Z.kEuZa6KJL.Xj8R1XhYYqgwqXhvZzEXpYf.', 'staff_warorot@rota.local', '22222222-2222-2222-2222-222222222222');