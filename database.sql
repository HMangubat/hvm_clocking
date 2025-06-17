-- Drop existing tables (for dev reset)
DROP TABLE IF EXISTS AuditLogs, Clockings, RaceResults, RaceParticipants, Races, Devices, LoftCoordinates, Pigeons, Users, Clubs CASCADE;

-- ========== USERS ==========
CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    full_name VARCHAR(100),
    email VARCHAR(100),
    phone_number VARCHAR(20),
    role VARCHAR(20) DEFAULT 'racer',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ========== CLUBS ==========
CREATE TABLE Clubs (
    club_id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    location VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ========== DEVICES ==========
CREATE TABLE Devices (
    device_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(user_id),
    name VARCHAR(100),
    serial_number VARCHAR(100) UNIQUE,
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ========== LOFT ==========
CREATE TABLE LoftCoordinates (
    loft_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(user_id) UNIQUE ON DELETE CASCADE,
    latitude DECIMAL(9,6) NOT NULL,
    longitude DECIMAL(9,6) NOT NULL
);

-- ========== PIGEONS ==========
CREATE TABLE Pigeons (
    pigeon_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    ring_number VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100),
    color VARCHAR(50),
    sex VARCHAR(10),
    breed VARCHAR(50),
    birth_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ========== RACES ==========
CREATE TABLE Races (
    race_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(100),
    distance_km DECIMAL(10, 2),
    release_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ========== PARTICIPANTS ==========
CREATE TABLE RaceParticipants (
    id SERIAL PRIMARY KEY,
    race_id INT REFERENCES Races(race_id) ON DELETE CASCADE,
    pigeon_id INT REFERENCES Pigeons(pigeon_id) ON DELETE CASCADE,
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (race_id, pigeon_id)
);

-- ========== CLOCKINGS ==========
CREATE TABLE Clockings (
    clocking_id SERIAL PRIMARY KEY,
    pigeon_id INT REFERENCES Pigeons(pigeon_id) ON DELETE CASCADE,
    race_id INT REFERENCES Races(race_id) ON DELETE CASCADE,
    user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    device_id INT,
    arrival_time TIMESTAMP NOT NULL,
    speed_kph DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ========== RACE RESULTS ==========
CREATE TABLE RaceResults (
    id SERIAL PRIMARY KEY,
    race_id INT REFERENCES Races(race_id) ON DELETE CASCADE,
    pigeon_id INT REFERENCES Pigeons(pigeon_id),
    speed_kph DECIMAL(10,2),
    arrival_time TIMESTAMP,
    rank INT
);

-- ========== AUDIT LOGS ==========
CREATE TABLE AuditLogs (
    log_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(user_id),
    action TEXT,
    log_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ========== SEED DATA ==========

-- Clubs
INSERT INTO Clubs (name, location) VALUES
('Sky Flyers Club', 'Manila'),
('Northwind Racers', 'Quezon City');

-- Users
INSERT INTO Users (username, password_hash, full_name, email, phone_number, role)
VALUES
('evcauyan', '$2a$10$zFZlKc5A7QeY8HxUwTe68.wGjMoVXJTxM1gZAZ6FKzEX3I0Rj5myy', 'Eric Cauyan', 'eric@example.com', '09171234567', 'admin'), -- password: 123456
('jmendoza', '$2a$10$zFZlKc5A7QeY8HxUwTe68.wGjMoVXJTxM1gZAZ6FKzEX3I0Rj5myy', 'Juan Mendoza', 'juan@example.com', '09181234567', 'racer');

-- Loft Coordinates
INSERT INTO LoftCoordinates (user_id, latitude, longitude)
VALUES
(1, 14.5995, 120.9842),
(2, 14.6760, 121.0437);

-- Devices
INSERT INTO Devices (user_id, name, serial_number) VALUES
(1, 'ClockMaster X100', 'DEV10001'),
(2, 'SpeedTracker Z200', 'DEV20001');

-- Pigeons
INSERT INTO Pigeons (user_id, ring_number, name, color, sex, breed, birth_date)
VALUES
(1, 'PH2024-001', 'Storm', 'Blue Bar', 'Male', 'Belgian', '2024-01-05'),
(1, 'PH2024-002', 'Shadow', 'Black', 'Female', 'German', '2024-02-10'),
(2, 'PH2024-003', 'Windchaser', 'White', 'Male', 'Dutch', '2024-03-15');

-- Races
INSERT INTO Races (name, location, distance_km, release_time)
VALUES
('Opening Race', 'Bulacan', 50.0, '2025-06-10 06:00:00'),
('Speed Derby', 'Pampanga', 100.0, '2025-06-12 06:00:00');

-- Participants
INSERT INTO RaceParticipants (race_id, pigeon_id) VALUES
(1, 1),
(1, 2),
(2, 3);

-- Clockings
INSERT INTO Clockings (pigeon_id, race_id, user_id, device_id, arrival_time, speed_kph)
VALUES
(1, 1, 1, 1, '2025-06-10 06:35:00', 85.71),
(2, 1, 1, 1, '2025-06-10 06:40:00', 75.00),
(3, 2, 2, 2, '2025-06-12 07:10:00', 84.00);

-- RaceResults
INSERT INTO RaceResults (race_id, pigeon_id, speed_kph, arrival_time, rank)
VALUES
(1, 1, 85.71, '2025-06-10 06:35:00', 1),
(1, 2, 75.00, '2025-06-10 06:40:00', 2),
(2, 3, 84.00, '2025-06-12 07:10:00', 1);

-- Audit Logs
INSERT INTO AuditLogs (user_id, action)
VALUES
(1, 'Registered pigeon PH2024-001'),
(2, 'Joined race Speed Derby with pigeon PH2024-003');
