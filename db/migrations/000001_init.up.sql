CREATE TABLE IF NOT EXISTS evaluations (
                                           id SERIAL PRIMARY KEY,
                                           expression TEXT NOT NULL
);


CREATE TABLE IF NOT EXISTS prime_evaluations (
    id SERIAL PRIMARY KEY,
    parent_id INT NOT NULL REFERENCES evaluations(id) ON DELETE CASCADE,
    arg1 DOUBLE PRECISION NOT NULL,
    arg2 DOUBLE PRECISION NOT NULL,
    operation TEXT NOT NULL,
    operation_time INT NOT NULL,
    result DOUBLE PRECISION NOT NULL,
    error BOOLEAN DEFAULT FALSE,
    completed_at timestamp with time zone
);


CREATE TABLE users (
                       id UUID PRIMARY KEY,
                       email TEXT UNIQUE NOT NULL,
                       hash TEXT NOT NULL,
                       created_at TIMESTAMP NOT NULL
);
