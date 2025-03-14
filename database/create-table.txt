-- SQLite
-- สร้างตาราง Users
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE,
    password TEXT,
    email TEXT UNIQUE NOT NULL,
    provider TEXT DEFAULT 'local',
    provider_id TEXT,
    profile_picture TEXT,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- สร้างตาราง OAuth Tokens
CREATE TABLE oauth_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    provider TEXT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    expires_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- สร้างตาราง Favorites
CREATE TABLE favorites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    station_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (station_id) REFERENCES stations(id)
);

-- สร้างตาราง Stations
CREATE TABLE stations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    location TEXT
);

-- สร้างตาราง Routes
CREATE TABLE routes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    start_station_id INTEGER,
    end_station_id INTEGER,
    distance REAL NOT NULL,
    duration TEXT NOT NULL,
    FOREIGN KEY (start_station_id) REFERENCES stations(id),
    FOREIGN KEY (end_station_id) REFERENCES stations(id)
);

-- สร้างตาราง Schedules
CREATE TABLE schedules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    route_id INTEGER,
    station_id INTEGER,
    round INTEGER NOT NULL,
    departure_time TIME NOT NULL,
    arrival_time TIME NOT NULL,
    FOREIGN KEY (route_id) REFERENCES routes(id),
    FOREIGN KEY (station_id) REFERENCES stations(id)
);

-- สร้างตาราง Schedule Logs
CREATE TABLE schedule_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    schedule_id INTEGER,
    staff_id INTEGER,
    change_description TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (schedule_id) REFERENCES schedules(id),
    FOREIGN KEY (staff_id) REFERENCES staff(id)
);

-- สร้างตาราง Staff
CREATE TABLE staff (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    station_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (station_id) REFERENCES stations(id)
);

-- สร้างตาราง Vehicles
CREATE TABLE vehicles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    license_plate TEXT UNIQUE NOT NULL,
    capacity INTEGER NOT NULL,
    driver_name TEXT NOT NULL,
    route_id INTEGER,
    FOREIGN KEY (route_id) REFERENCES routes(id)
);
