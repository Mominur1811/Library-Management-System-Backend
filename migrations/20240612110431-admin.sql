
-- +migrate Up
CREATE TABLE IF NOT EXISTs admin(
    Id SERIAL PRIMARY KEY,
    Email VARCHAR(40) UNIQUE,
    Password VARCHAR(80),
    Is_SuperAdmin BOOLEAN DEFAULT false
);

-- +migrate Down
DROP TABLE IF EXISTs admin;