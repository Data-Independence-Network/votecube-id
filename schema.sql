CREATE DATABASE votecube;

CREATE TABLE votecube.users(
    id INT PRIMARY KEY, 
    username STRING NOT NULL, 
    password STRING,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
