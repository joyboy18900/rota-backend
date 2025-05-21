-- Remove seeded data in reverse order
DELETE FROM staff WHERE id IN (
    'staff1111-1111-1111-1111-111111111111',
    'staff2222-2222-2222-2222-222222222222'
);

DELETE FROM users WHERE id = 'admin1111-1111-1111-1111-111111111111';

DELETE FROM route_vehicles;
DELETE FROM route_stations;
DELETE FROM vehicles;
DELETE FROM routes;
DELETE FROM stations;