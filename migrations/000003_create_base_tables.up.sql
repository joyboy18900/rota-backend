-- Create stations table
CREATE TABLE IF NOT EXISTS stations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location TEXT NOT NULL,
    detail TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create routes table
CREATE TABLE IF NOT EXISTS routes (
    id SERIAL PRIMARY KEY,
    start_station_id INTEGER REFERENCES stations(id) NOT NULL,
    end_station_id INTEGER REFERENCES stations(id) NOT NULL,
    distance FLOAT NOT NULL,
    duration VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create vehicles table
CREATE TABLE IF NOT EXISTS vehicles (
    id SERIAL PRIMARY KEY,
    license_plate VARCHAR(20) UNIQUE NOT NULL,
    capacity INTEGER NOT NULL,
    driver_name VARCHAR(100),
    route_id INTEGER REFERENCES routes(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE,
    password TEXT,
    email VARCHAR(100) UNIQUE NOT NULL,
    provider VARCHAR(20) NOT NULL DEFAULT 'local',
    provider_id TEXT,
    profile_picture TEXT,
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create oauth_tokens table
CREATE TABLE IF NOT EXISTS oauth_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    provider VARCHAR(20) NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create staff table
CREATE TABLE IF NOT EXISTS staff (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    station_id INTEGER REFERENCES stations(id) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create schedules table
CREATE TABLE IF NOT EXISTS schedules (
    id SERIAL PRIMARY KEY,
    route_id INTEGER REFERENCES routes(id) NOT NULL,
    station_id INTEGER REFERENCES stations(id) NOT NULL,
    round INTEGER NOT NULL CHECK (round BETWEEN 1 AND 20),
    departure_time TIME NOT NULL,
    arrival_time TIME NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create schedule_logs table
CREATE TABLE IF NOT EXISTS schedule_logs (
    id SERIAL PRIMARY KEY,
    schedule_id INTEGER REFERENCES schedules(id) NOT NULL,
    staff_id INTEGER REFERENCES staff(id) NOT NULL,
    change_description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create favorites table
CREATE TABLE IF NOT EXISTS favorites (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    station_id INTEGER REFERENCES stations(id) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, station_id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_stations_deleted_at ON stations(deleted_at);
CREATE INDEX IF NOT EXISTS idx_routes_deleted_at ON routes(deleted_at);
CREATE INDEX IF NOT EXISTS idx_vehicles_deleted_at ON vehicles(deleted_at);
CREATE INDEX IF NOT EXISTS idx_staff_deleted_at ON staff(deleted_at);
CREATE INDEX IF NOT EXISTS idx_schedules_deleted_at ON schedules(deleted_at);
CREATE INDEX IF NOT EXISTS idx_route_stations_route_id ON route_stations(route_id);
CREATE INDEX IF NOT EXISTS idx_route_stations_station_id ON route_stations(station_id);
CREATE INDEX IF NOT EXISTS idx_route_vehicles_route_id ON route_vehicles(route_id);
CREATE INDEX IF NOT EXISTS idx_route_vehicles_vehicle_id ON route_vehicles(vehicle_id);
CREATE INDEX IF NOT EXISTS idx_schedule_logs_schedule_id ON schedule_logs(schedule_id);
CREATE INDEX IF NOT EXISTS idx_schedule_logs_staff_id ON schedule_logs(staff_id);
CREATE INDEX IF NOT EXISTS idx_favorites_user_id ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_station_id ON favorites(station_id);