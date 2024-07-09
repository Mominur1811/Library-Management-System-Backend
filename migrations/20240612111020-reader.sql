
-- +migrate Up
CREATE TABLE IF NOT EXISTs reader(
    Id SERIAL PRIMARY KEY,
    Name VARCHAR,
    Email VARCHAR UNIQUE,
    Password VARCHAR,
    Is_Active BOOLEAN DEFAULT false
);

-- +migrate Down
DROP TABLE IF EXISTS reader;