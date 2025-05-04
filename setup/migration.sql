create database postgres;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uuid varchar(100) not null,
    username VARCHAR(255) NOT null unique,
    password VARCHAR(255) NOT null,
    email varchar(255) not null,
    created_at timestamp not null
);