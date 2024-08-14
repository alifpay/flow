CREATE TABLE rules(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    conditions jsonb NOT NULL
);

CREATE TABLE flows(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    
);