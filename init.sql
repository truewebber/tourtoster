DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS hotel;

CREATE TABLE users
(
    id            INTEGER PRIMARY KEY autoincrement,
    first_name    VARCHAR(60) NOT NULL,
    second_name   VARCHAR(60) NOT NULL,
    last_name     VARCHAR(60) NOT NULL,
    hotel_name    VARCHAR(100),
    hotel_id      INTEGER NOT NULL DEFAULT 0,
    note          TEXT,
    email         VARCHAR(60)        NOT NULL,
    phone         VARCHAR(20)        NOT NULL,
    password_hash TEXT,
    status        TINYINT   DEFAULT 1,
    role          INTEGER,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX idx_users_email ON users (email);
CREATE UNIQUE INDEX idx_users_phone ON users (phone);

insert into users (first_name, second_name, last_name, hotel_name, note, email, phone, password_hash, role)
VALUES ('', '', '', '', '', 'kish94@mail.ru', '+79643896032', '', 0);

CREATE TABLE hotel
(
    id   INTEGER PRIMARY KEY autoincrement,
    name VARCHAR(100)
);
CREATE UNIQUE INDEX idx_hotel_name ON hotel (name);