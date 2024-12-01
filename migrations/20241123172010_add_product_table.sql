-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(100) NOT NULL,
                          description TEXT,
                          price NUMERIC(10, 2) NOT NULL,
                          user_id INT REFERENCES users(id),
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
DROP TABLE products;
-- +goose StatementEnd
-- +goose Down

