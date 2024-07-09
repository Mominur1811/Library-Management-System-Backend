
-- +migrate Up
CREATE TABLE IF NOT EXISTS book (
    Id SERIAL PRIMARY KEY,
    Title VARCHAR NOT NULL,
    Category VARCHAR NOT NULL,
    Author VARCHAR NOT NULL,
    Summary VARCHAR NOT NULL,
    Quantity INT NOT NULL,
    Available INT NOT NULL,
    Total_page INT NOT NULL,
    Image_link  VARCHAR NOT NULL,
    Created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +migrate Down
DROP TABLE IF EXISTS book;