CREATE TABLE IF NOT EXISTS urls(
    id serial PRIMARY KEY,
    original VARCHAR(255),
    touched TIMESTAMP
);