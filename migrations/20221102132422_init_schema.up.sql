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

CREATE INDEX IF NOT EXISTS idx_time_touched
    ON urls (touched);

INSERT INTO count (id,value)
VALUES(true,0)
ON CONFLICT (id)
    DO NOTHING;

