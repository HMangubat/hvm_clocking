-- Include your full SQL schema from earlier here
-- USERS
CREATE TABLE Users (
    user_id         SERIAL PRIMARY KEY,
    username        VARCHAR(50) UNIQUE NOT NULL,
    password_hash   TEXT NOT NULL,
    full_name       VARCHAR(100),
    email           VARCHAR(100),
    phone_number    VARCHAR(20),
    role            VARCHAR(20) DEFAULT 'fancier',
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add other table definitions...
-- (Lofts, Pigeons, ClockingDevices, Races, RaceEntries, Clockings, RaceResults)
