CREATE TABLE flows(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    node jsonb NOT NULL
);

CREATE TABLE data_inputs(
    id VARCHAR(63) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
);

CREATE TABLE functions(
    id VARCHAR(63) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
);