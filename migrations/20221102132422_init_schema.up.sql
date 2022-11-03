CREATE TABLE IF NOT EXISTS urls(
    id serial PRIMARY KEY,
    original VARCHAR(255) UNIQUE,
    touched TIMESTAMP
);

CREATE TABLE IF NOT EXISTS count(
    id bool PRIMARY KEY DEFAULT TRUE,
    value INT NOT NULL,
    CONSTRAINT count_uni CHECK (id)
);

INSERT INTO count (value) VALUES(0);
