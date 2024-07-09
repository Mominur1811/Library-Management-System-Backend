
-- +migrate Up
CREATE TABLE IF NOT EXISTS borrow_history (
    request_id SERIAL PRIMARY KEY,
    book_id INT,
    borrower_id INT,
    issued_at TIMESTAMP,
    returned_at TIMESTAMP,
    borrow_status VARCHAR,
    read_page INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +migrate Down
DROP TABLE IF EXISTS borrow_history;
