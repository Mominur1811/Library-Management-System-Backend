
-- +migrate Up
CREATE TABLE IF NOT EXISTS book_request (
    Request_id SERIAL PRIMARY KEY,
    Bookid INT,
    Readerid INT,
    Issued_at TIMESTAMP,
    Request_status VARCHAR,
    read_page INT DEFAULT 0,
    Created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS book_request;