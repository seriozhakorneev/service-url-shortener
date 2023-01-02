CREATE TABLE IF NOT EXISTS urls(
    id serial PRIMARY KEY,
    original VARCHAR(255) UNIQUE,
    live_till TIMESTAMP,
    activated TIMESTAMP
);

CREATE TABLE IF NOT EXISTS last(
    id bool PRIMARY KEY DEFAULT TRUE,
    value INT NOT NULL,
    CONSTRAINT last_uni CHECK (id)
);

CREATE INDEX IF NOT EXISTS idx_time_activated
    ON urls (activated);

INSERT INTO last (id,value)
VALUES(true,0)
ON CONFLICT (id)
    DO NOTHING;

