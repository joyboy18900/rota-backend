CREATE TABLE IF NOT EXISTS staffs (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255),
    password VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    name VARCHAR(255) NOT NULL,
    position VARCHAR(255),
    phone VARCHAR(50),
    station_id INTEGER REFERENCES stations(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
