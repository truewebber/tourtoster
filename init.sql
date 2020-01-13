CREATE TABLE users
(
    id            INTEGER PRIMARY KEY autoincrement,
    first_name    VARCHAR(60) NOT NULL,
    second_name   VARCHAR(60) NOT NULL,
    last_name     VARCHAR(60) NOT NULL,
    hotel_name    VARCHAR(60),
    hotel_id      INTEGER,
    note          TEXT,
    email         TEXT        NOT NULL UNIQUE,
    phone         TEXT        NOT NULL UNIQUE,
    password_hash TEXT,
    status        TINYINT   DEFAULT 1,
    role          INTEGER,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- insert into users (first_name, second_name, last_name, hotel_name, hotel_id, note, email, phone, password_hash, role)
-- VALUES ('', '', '', '', 0, '', 'kish94@mail.ru', '+79643896032', '', 0);
