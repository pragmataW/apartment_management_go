CREATE DATABASE apartments;

\c apartments;

CREATE TABLE apartments (
    flat_no INT PRIMARY KEY,
    owner_name VARCHAR(255),
    owner_surname VARCHAR(255), 
    mail VARCHAR(255) UNIQUE,           
    password VARCHAR(255),
    dues_count INT
);

CREATE TABLE announcements (
    announcement_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL
);

CREATE TABLE merchants (
    merchant_id VARCHAR(255) PRIMARY KEY,
    email TEXT NOT NULL
);