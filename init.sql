CREATE TABLE IF NOT EXISTS users
(
    id            INTEGER PRIMARY KEY autoincrement,
    first_name    VARCHAR(60) NOT NULL,
    second_name   VARCHAR(60) NOT NULL,
    last_name     VARCHAR(60) NOT NULL,
    hotel_name    VARCHAR(100),
    hotel_id      INTEGER     NOT NULL DEFAULT 0,
    note          TEXT,
    email         VARCHAR(60) NOT NULL,
    phone         VARCHAR(20) NOT NULL,
    password_hash TEXT,
    status        TINYINT              DEFAULT 1,
    role          INTEGER,
    updated_at    TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    created_at    TIMESTAMP            DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS idx__users__email ON users (email);
CREATE UNIQUE INDEX IF NOT EXISTS idx__users__phone ON users (phone);

INSERT INTO users (first_name, second_name, last_name, hotel_name, note, email, phone, password_hash, role)
VALUES ('Aleksey', '', 'Kish', 'Blahotel', '', 'kish94@mail.ru', '+79643896032', '', 0)
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS hotel
(
    id   INTEGER PRIMARY KEY autoincrement,
    name VARCHAR(100)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx__hotel__name ON hotel (name);

-- ###########################################################################################

PRAGMA foreign_keys=1;

CREATE TABLE IF NOT EXISTS tour_types
(
    id   INTEGER PRIMARY KEY autoincrement,
    name VARCHAR NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx__tour_types__name ON tour_types (name);
INSERT INTO tour_types (name)
VALUES ('group tour'),
       ('private tour')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS tours
(
    id              INTEGER PRIMARY KEY autoincrement,
    tour_type_id    SMALLINT NOT NULL,
    creator_id      INTEGER  NOT NULL,
    title           VARCHAR  NOT NULL,
    image           TEXT      DEFAULT NULL,
    description     TEXT      DEFAULT NULL,
    map_description TEXT      DEFAULT NULL,
    max_persons     SMALLINT NOT NULL,
    price_per_adult INTEGER  NOT NULL,
    price_per_child INTEGER  NOT NULL,
    status          SMALLINT NOT NULL,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tour_type_id) REFERENCES tour_types (id)
);
CREATE INDEX IF NOT EXISTS idx__tours__creator_id ON tours (creator_id);

CREATE TABLE IF NOT EXISTS features
(
    id           INTEGER PRIMARY KEY autoincrement,
    tour_type_id SMALLINT NOT NULL,
    icon         TEXT     NOT NULL,
    title        TEXT     NOT NULL,
    FOREIGN KEY (tour_type_id) REFERENCES tour_types (id)
);
INSERT INTO features (tour_type_id, title, icon)
VALUES (1, 'AVAILABILITY',
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 56"><defs><style>.cls-1{fill:#fff;}.cls-2{fill:#4d5152;}.cls-3{fill:#67dde0;}.cls-4{fill:#f28c13;}</style></defs><g id="Layer_2" data-name="Layer 2"><g id="Layer_1-2" data-name="Layer 1"><rect class="cls-1" x="1" y="11" width="62" height="44"/><rect class="cls-1" x="1" y="1" width="62" height="10"/><path class="cls-2" d="M63,56H1a1,1,0,0,1-1-1V11a1,1,0,0,1,1-1H63a1,1,0,0,1,1,1V55A1,1,0,0,1,63,56ZM2,54H62V12H2Z"/><rect class="cls-3" x="21" y="18" width="6" height="6"/><rect class="cls-3" x="29" y="18" width="6" height="6"/><rect class="cls-3" x="37" y="18" width="6" height="6"/><rect class="cls-4" x="45" y="18" width="6" height="6"/><rect class="cls-4" x="53" y="18" width="6" height="6"/><rect class="cls-3" x="5" y="26" width="6" height="6"/><rect class="cls-3" x="13" y="26" width="6" height="6"/><rect class="cls-3" x="21" y="26" width="6" height="6"/><rect class="cls-3" x="29" y="26" width="6" height="6"/><rect class="cls-3" x="37" y="26" width="6" height="6"/><rect class="cls-4" x="45" y="26" width="6" height="6"/><rect class="cls-4" x="53" y="26" width="6" height="6"/><rect class="cls-3" x="5" y="34" width="6" height="6"/><rect class="cls-3" x="13" y="34" width="6" height="6"/><rect class="cls-3" x="21" y="34" width="6" height="6"/><rect class="cls-3" x="29" y="34" width="6" height="6"/><rect class="cls-3" x="37" y="34" width="6" height="6"/><rect class="cls-4" x="45" y="34" width="6" height="6"/><rect class="cls-4" x="53" y="34" width="6" height="6"/><rect class="cls-3" x="5" y="42" width="6" height="6"/><rect class="cls-3" x="13" y="42" width="6" height="6"/><rect class="cls-3" x="21" y="42" width="6" height="6"/><rect class="cls-3" x="29" y="42" width="6" height="6"/><rect class="cls-3" x="37" y="42" width="6" height="6"/><path class="cls-2" d="M64,8H62V2H2V8H0V1A1,1,0,0,1,1,0H63a1,1,0,0,1,1,1Z"/><rect class="cls-2" x="7" y="4" width="50" height="2"/><rect class="cls-2" x="35.76" y="20" width="8.49" height="2" transform="translate(-3.13 34.44) rotate(-45)"/><rect class="cls-2" x="39" y="16.76" width="2" height="8.49" transform="translate(-3.13 34.44) rotate(-45)"/><rect class="cls-2" x="11.76" y="36" width="8.49" height="2" transform="translate(-21.48 22.15) rotate(-45)"/><rect class="cls-2" x="15" y="32.76" width="2" height="8.49" transform="translate(-21.48 22.15) rotate(-45)"/><rect class="cls-2" x="15" y="20" width="2" height="2"/><rect class="cls-2" x="7" y="20" width="2" height="2"/><rect class="cls-2" x="55" y="44" width="2" height="2"/><rect class="cls-2" x="47" y="44" width="2" height="2"/></g></g></svg>'),
       (1, 'DURATION',
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 56 64"><defs><style>.cls-1{fill:#fff;}.cls-2{fill:#4d5152;}.cls-3{fill:#67dde0;}.cls-4{fill:#f28c13;}</style></defs><g id="Layer_2" data-name="Layer 2"><g id="Layer_1-2" data-name="Layer 1"><circle class="cls-1" cx="28" cy="36" r="27"/><path class="cls-2" d="M28,64A28,28,0,1,1,56,36,28,28,0,0,1,28,64Zm0-54A26,26,0,1,0,54,36,26,26,0,0,0,28,10Z"/><path class="cls-2" d="M28,10a5,5,0,1,1,5-5A5,5,0,0,1,28,10Zm0-8a3,3,0,1,0,3,3A3,3,0,0,0,28,2Z"/><rect class="cls-2" x="27" y="5" width="2" height="7"/><circle class="cls-3" cx="28" cy="36" r="20"/><path class="cls-4" d="M28,16V36H48A20,20,0,0,0,28,16Z"/><circle class="cls-1" cx="28" cy="36" r="6"/><rect class="cls-2" x="28" y="35" width="20" height="2"/><rect class="cls-2" x="46.38" y="14.5" width="4.24" height="2" transform="translate(3.25 38.83) rotate(-45)"/><rect class="cls-2" x="6.5" y="13.38" width="2" height="4.25" transform="translate(-8.76 9.84) rotate(-45)"/><rect class="cls-2" x="49" y="11.17" width="2" height="5.66" transform="translate(4.75 39.46) rotate(-45)"/><rect class="cls-2" x="3.17" y="13" width="5.66" height="2" transform="translate(-8.14 8.34) rotate(-45)"/><circle class="cls-2" cx="28" cy="36" r="2"/><rect class="cls-2" x="1" y="35" width="3" height="2"/><rect class="cls-2" x="52" y="35" width="3" height="2"/><rect class="cls-2" x="27" y="60" width="2" height="3"/></g></g></svg>'),
       (1, 'MEET ON LOCATION',
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 64"><defs><style>.cls-1{fill:#67dde0;}.cls-2{fill:#acf0f2;}.cls-3{fill:#ffb957;}.cls-4{fill:#4d5152;}.cls-5{fill:#fff;}</style></defs><g id="Layer_2" data-name="Layer 2"><g id="Layer_1-2" data-name="Layer 1"><ellipse class="cls-1" cx="24" cy="54" rx="23" ry="10"/><ellipse class="cls-2" cx="24" cy="53" rx="9" ry="4"/><path class="cls-3" d="M43,20c0,14-19,33-19,33S5,34,5,20a19,19,0,0,1,38,0Z"/><path class="cls-4" d="M24,54.41l-.71-.7C22.5,52.92,4,34.25,4,20a20,20,0,0,1,40,0c0,14.25-18.5,32.92-19.29,33.71ZM24,2A18,18,0,0,0,6,20C6,31.92,20.59,48,24,51.56,27.41,48,42,31.91,42,20A18,18,0,0,0,24,2Z"/><circle class="cls-5" cx="24" cy="20" r="12"/><path class="cls-4" d="M24,28a8,8,0,1,1,8-8A8,8,0,0,1,24,28Zm0-14a6,6,0,1,0,6,6A6,6,0,0,0,24,14Z"/><path class="cls-4" d="M24,64C10.54,64,0,59.17,0,53c0-3.59,3.53-6.83,9.68-8.88L10.32,46C5.11,47.75,2,50.37,2,53c0,4.88,10.08,9,22,9s22-4.12,22-9c0-2.64-3.11-5.25-8.32-7l.63-1.89C44.47,46.17,48,49.41,48,53,48,59.17,37.46,64,24,64Z"/><rect class="cls-4" x="23" y="19" width="2" height="2"/></g></g></svg>'),
       (1, 'LANGUAGE',
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 60"><defs><style>.cls-1{fill:#fff;}.cls-2{fill:#67dde0;}.cls-3{fill:#4d5152;}.cls-4{fill:#ffb957;}</style></defs><g id="Layer_2" data-name="Layer 2"><g id="Layer_1-2" data-name="Layer 1"><polygon class="cls-1" points="42 25 63 25 63 1 23 1 23 14 23 31 36 31 42 25"/><path class="cls-2" d="M31,44c7,4,14,3,14,13v2H1V57c0-10,7-9,14-13Z"/><polygon class="cls-1" points="27 48 23 46 19 48 15 44 17 41 18 41 18 37 23 38 28 37 28 41 29 41 31 44 27 48"/><path class="cls-3" d="M46,60H0V57c0-8,4.45-9.7,9.15-11.49a35.44,35.44,0,0,0,5.35-2.38l1,1.74a39.47,39.47,0,0,1-5.64,2.51C5.13,49.18,2,50.37,2,57v1H44V57c0-6.63-3.13-7.82-7.86-9.62a39.47,39.47,0,0,1-5.64-2.51l1-1.74a35.44,35.44,0,0,0,5.35,2.38C41.55,47.3,46,49,46,57Z"/><path class="cls-3" d="M27.2,49.22,23,47.12l-4.2,2.1-5.09-5.09,2.82-4.22L23,41l6.47-1.08,2.82,4.22ZM23,44.88l3.8,1.9,2.91-2.91-1.18-1.78L23,43l-5.53-.92-1.18,1.78,2.91,2.91Z"/><rect class="cls-3" x="27" y="37" width="2" height="4"/><rect class="cls-3" x="17" y="37" width="2" height="4"/><rect class="cls-3" x="22" y="54" width="2" height="2"/><rect class="cls-3" x="22" y="50" width="2" height="2"/><path class="cls-1" d="M32,24v4c-1,7-4,9-4,9a13,13,0,0,1-10,0s-3-2-4-9V24l4-2,10,4Z"/><path class="cls-4" d="M14,24c0-6.08,2.92-9,9-9s9,2.92,9,9l-4,2-5-2-5-2Z"/><path class="cls-3" d="M23,39a14.21,14.21,0,0,1-5.45-1.11l-.1-.06c-.14-.09-3.39-2.33-4.44-9.69V24h2v3.93c.86,5.92,3.2,8,3.52,8.21a12.1,12.1,0,0,0,9,0c.31-.25,2.66-2.3,3.52-8.21V24h2v4.14c-1.05,7.36-4.3,9.6-4.44,9.69l-.1.06A14.21,14.21,0,0,1,23,39Z"/><path class="cls-3" d="M33,24H31c0-5.53-2.47-8-8-8s-8,2.47-8,8H13c0-6.64,3.36-10,10-10S33,17.36,33,24Z"/><polygon class="cls-3" points="28.05 27.1 18.05 23.1 14.45 24.89 13.55 23.11 17.95 20.9 27.95 24.9 31.55 23.11 32.45 24.89 28.05 27.1"/><polygon class="cls-3" points="64 26 42 26 42 24 62 24 62 2 24 2 24 12 22 12 22 0 64 0 64 26"/><path class="cls-3" d="M36,32a1,1,0,0,1-.71-.29,1,1,0,0,1,0-1.42l6-6a1,1,0,0,1,1.42,1.42l-6,6A1,1,0,0,1,36,32Z"/><rect class="cls-3" x="38" y="12" width="2" height="2"/><rect class="cls-3" x="42" y="12" width="2" height="2"/><rect class="cls-3" x="46" y="12" width="2" height="2"/></g></g></svg>'),
       (1, 'GROUP SIZE',
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 46"><defs><style>.cls-1{fill:#acf0f2;}.cls-2{fill:#ffb957;}.cls-3{fill:#4d5152;}.cls-4{fill:#67dde0;}.cls-5{fill:#fff;}</style></defs><g id="Layer_2" data-name="Layer 2"><g id="Layer_1-2" data-name="Layer 1"><path class="cls-1" d="M55,27a21.77,21.77,0,0,0,5,3c3,1,3,5,3,5v2H41V35s0-4,3-5a21.77,21.77,0,0,0,5-3Z"/><path class="cls-2" d="M57,20.5C57,24,55,27,52,27s-5-3-5-6.5S48,15,52,15,57,17,57,20.5Z"/><path class="cls-3" d="M52,28c-3.42,0-6-3.22-6-7.5,0-3,.68-6.5,6-6.5s6,3.47,6,6.5C58,24.78,55.42,28,52,28Zm0-12c-3,0-4,1.09-4,4.5,0,2.67,1.4,5.5,4,5.5s4-2.83,4-5.5C56,17.09,55,16,52,16Z"/><path class="cls-3" d="M64,38H40V35c0-.19,0-4.73,3.68-5.95a20.86,20.86,0,0,0,4.7-2.83l1.24,1.56A22.34,22.34,0,0,1,44.32,31C42,31.71,42,35,42,35v1H62V35s0-3.29-2.32-4.05a22.34,22.34,0,0,1-5.3-3.17l1.24-1.56a20.86,20.86,0,0,0,4.7,2.83C64,30.27,64,34.81,64,35Z"/><path class="cls-1" d="M15,27a21.77,21.77,0,0,0,5,3c3,1,3,5,3,5v2H1V35s0-4,3-5a21.77,21.77,0,0,0,5-3Z"/><path class="cls-2" d="M17,20.5C17,24,15,27,12,27s-5-3-5-6.5S8,15,12,15,17,17,17,20.5Z"/><path class="cls-3" d="M12,28c-3.42,0-6-3.22-6-7.5,0-3,.68-6.5,6-6.5s6,3.47,6,6.5C18,24.78,15.42,28,12,28Zm0-12c-3,0-4,1.09-4,4.5C8,23.17,9.4,26,12,26s4-2.83,4-5.5C16,17.09,15,16,12,16Z"/><path class="cls-3" d="M24,38H0V35c0-.19,0-4.73,3.68-5.95a20.86,20.86,0,0,0,4.7-2.83l1.24,1.56A22.34,22.34,0,0,1,4.32,31C2,31.71,2,35,2,35v1H22V35S22,31.71,19.68,31a22.34,22.34,0,0,1-5.3-3.17l1.24-1.56a20.86,20.86,0,0,0,4.7,2.83C24,30.27,24,34.81,24,35Z"/><path class="cls-4" d="M40,30c7,4,14,3,14,13v2H10V43c0-10,7-9,14-13Z"/><polygon class="cls-5" points="36 34 32 32 28 34 24 30 26 27 27 27 27 23 32 24 37 23 37 27 38 27 40 30 36 34"/><path class="cls-3" d="M55,46H9V43c0-8,4.45-9.7,9.15-11.49a35.44,35.44,0,0,0,5.35-2.38l1,1.74a39.47,39.47,0,0,1-5.64,2.51C14.13,35.18,11,36.37,11,43v1H53V43c0-6.63-3.13-7.82-7.86-9.62a39.47,39.47,0,0,1-5.64-2.51l1-1.74a35.44,35.44,0,0,0,5.35,2.38C50.55,33.3,55,35,55,43Z"/><path class="cls-3" d="M36.2,35.22,32,33.12l-4.2,2.1-5.09-5.09,2.82-4.22L32,27l6.47-1.08,2.82,4.22ZM32,30.88l3.8,1.9,2.91-2.91-1.18-1.78L32,29l-5.53-.92-1.18,1.78,2.91,2.91Z"/><rect class="cls-3" x="36" y="23" width="2" height="4"/><rect class="cls-3" x="26" y="23" width="2" height="4"/><rect class="cls-3" x="31" y="40" width="2" height="2"/><rect class="cls-3" x="31" y="36" width="2" height="2"/><path class="cls-5" d="M41,10v4c-1,7-4,9-4,9a13,13,0,0,1-10,0s-3-2-4-9V10l4-2,10,4Z"/><path class="cls-2" d="M23,10c0-6.08,2.92-9,9-9s9,2.92,9,9l-4,2-5-2L27,8Z"/><path class="cls-3" d="M32,25a14.21,14.21,0,0,1-5.45-1.11l-.1-.06c-.14-.09-3.39-2.33-4.44-9.69V10h2v3.93c.86,5.92,3.2,8,3.52,8.21a12.1,12.1,0,0,0,9,0c.31-.25,2.66-2.3,3.52-8.21V10h2v4.14c-1.05,7.36-4.3,9.6-4.44,9.69l-.1.06A14.21,14.21,0,0,1,32,25Z"/><path class="cls-3" d="M42,10H40c0-5.53-2.47-8-8-8s-8,2.47-8,8H22C22,3.36,25.36,0,32,0S42,3.36,42,10Z"/><polygon class="cls-3" points="37.05 13.1 27.05 9.1 23.45 10.89 22.55 9.11 26.95 6.9 36.95 10.9 40.55 9.11 41.45 10.89 37.05 13.1"/></g></g></svg>'),
       (1, 'MIN. AGE',
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 56 60"><defs><style>.cls-1{fill:#acf0f2;}.cls-2{fill:#ffb957;}.cls-3{fill:#f28c13;}.cls-4{fill:#4d5152;}.cls-5{fill:#fff;}</style></defs><g id="Layer_2" data-name="Layer 2"><g id="Layer_1-2" data-name="Layer 1"><circle class="cls-1" cx="30" cy="30" r="26"/><path class="cls-2" d="M52,49,30,30,52,11h-.09a29,29,0,1,0,0,38Z"/><path class="cls-3" d="M9,24A28.89,28.89,0,0,1,16.88,4.14,29,29,0,1,0,51.13,49.86,29,29,0,0,1,9,24Z"/><path class="cls-4" d="M30,60A30,30,0,1,1,52.55,10.21a1,1,0,0,1,.35.45,1,1,0,0,1-.25,1.1L31.53,30,52.65,48.24a1,1,0,0,1,.29,1.11,1,1,0,0,1-.48.54A30,30,0,0,1,30,60ZM30,2A28,28,0,1,0,50.52,49.05L29.35,30.76a1,1,0,0,1,0-1.52L50.52,11A28,28,0,0,0,30,2Z"/><circle class="cls-5" cx="30" cy="14" r="4"/><rect class="cls-4" x="38" y="29" width="2" height="2"/><rect class="cls-4" x="44" y="29" width="2" height="2"/><rect class="cls-4" x="50" y="29" width="2" height="2"/></g></g></svg>'),
       (1, 'PICK-UP',
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 62 58"><defs><style>.cls-1{fill:#4d5152;}.cls-2{fill:#f28c13;}.cls-3{fill:#fff;}.cls-4{fill:#67dde0;}</style></defs><g id="Layer_2" data-name="Layer 2"><g id="Layer_1-2" data-name="Layer 1"><path class="cls-1" d="M61,58H29V56H60V4H28V54H26V3a1,1,0,0,1,1-1H61a1,1,0,0,1,1,1V57A1,1,0,0,1,61,58Z"/><rect class="cls-2" x="30" y="6" width="8" height="8"/><rect class="cls-2" x="40" y="6" width="8" height="8"/><rect class="cls-2" x="50" y="6" width="8" height="8"/><rect class="cls-2" x="40" y="16" width="8" height="8"/><rect class="cls-2" x="50" y="16" width="8" height="8"/><rect class="cls-2" x="40" y="26" width="8" height="8"/><rect class="cls-2" x="50" y="26" width="8" height="8"/><rect class="cls-2" x="40" y="36" width="8" height="8"/><rect class="cls-2" x="40" y="46" width="8" height="8"/><rect class="cls-2" x="50" y="36" width="8" height="8"/><rect class="cls-2" x="50" y="46" width="8" height="8"/><polygon class="cls-3" points="37 17 1 17 1 57 12 57 12 51 26 51 26 57 37 57 37 17"/><path class="cls-1" d="M37,58H28V56h8V18H2V56h8v2H1a1,1,0,0,1-1-1V17a1,1,0,0,1,1-1H37a1,1,0,0,1,1,1V57A1,1,0,0,1,37,58Z"/><rect class="cls-4" x="4" y="20" width="6" height="4"/><rect class="cls-4" x="12" y="20" width="6" height="4"/><rect class="cls-4" x="4" y="26" width="6" height="4"/><rect class="cls-4" x="12" y="26" width="6" height="4"/><rect class="cls-4" x="4" y="32" width="6" height="4"/><rect class="cls-4" x="12" y="32" width="6" height="4"/><rect class="cls-4" x="4" y="38" width="6" height="4"/><rect class="cls-4" x="12" y="38" width="6" height="4"/><rect class="cls-4" x="4" y="44" width="6" height="4"/><rect class="cls-4" x="12" y="44" width="6" height="4"/><rect class="cls-4" x="20" y="20" width="6" height="4"/><rect class="cls-4" x="28" y="20" width="6" height="4"/><rect class="cls-4" x="20" y="26" width="6" height="4"/><rect class="cls-4" x="28" y="26" width="6" height="4"/><rect class="cls-4" x="20" y="32" width="6" height="4"/><rect class="cls-4" x="28" y="32" width="6" height="4"/><rect class="cls-4" x="20" y="38" width="6" height="4"/><rect class="cls-4" x="28" y="38" width="6" height="4"/><rect class="cls-4" x="20" y="44" width="6" height="4"/><rect class="cls-4" x="28" y="44" width="6" height="4"/><rect class="cls-4" x="4" y="50" width="6" height="4"/><rect class="cls-4" x="28" y="50" width="6" height="4"/><rect class="cls-1" y="12" width="24" height="2"/><rect class="cls-1" y="8" width="24" height="2"/><rect class="cls-1" y="4" width="24" height="2"/><rect class="cls-1" width="24" height="2"/><path class="cls-1" d="M25,58H13a1,1,0,0,1-1-1V51a1,1,0,0,1,1-1H25a1,1,0,0,1,1,1v6A1,1,0,0,1,25,58ZM14,56H24V52H14Z"/></g></g></svg>'),
       (2, 'AVAILABILITY',
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 56"><defs><style>.cls-1{fill:#fff;}.cls-2{fill:#4d5152;}.cls-3{fill:#67dde0;}.cls-4{fill:#f28c13;}</style></defs><g id="Layer_2" data-name="Layer 2"><g id="Layer_1-2" data-name="Layer 1"><rect class="cls-1" x="1" y="11" width="62" height="44"/><rect class="cls-1" x="1" y="1" width="62" height="10"/><path class="cls-2" d="M63,56H1a1,1,0,0,1-1-1V11a1,1,0,0,1,1-1H63a1,1,0,0,1,1,1V55A1,1,0,0,1,63,56ZM2,54H62V12H2Z"/><rect class="cls-3" x="21" y="18" width="6" height="6"/><rect class="cls-3" x="29" y="18" width="6" height="6"/><rect class="cls-3" x="37" y="18" width="6" height="6"/><rect class="cls-4" x="45" y="18" width="6" height="6"/><rect class="cls-4" x="53" y="18" width="6" height="6"/><rect class="cls-3" x="5" y="26" width="6" height="6"/><rect class="cls-3" x="13" y="26" width="6" height="6"/><rect class="cls-3" x="21" y="26" width="6" height="6"/><rect class="cls-3" x="29" y="26" width="6" height="6"/><rect class="cls-3" x="37" y="26" width="6" height="6"/><rect class="cls-4" x="45" y="26" width="6" height="6"/><rect class="cls-4" x="53" y="26" width="6" height="6"/><rect class="cls-3" x="5" y="34" width="6" height="6"/><rect class="cls-3" x="13" y="34" width="6" height="6"/><rect class="cls-3" x="21" y="34" width="6" height="6"/><rect class="cls-3" x="29" y="34" width="6" height="6"/><rect class="cls-3" x="37" y="34" width="6" height="6"/><rect class="cls-4" x="45" y="34" width="6" height="6"/><rect class="cls-4" x="53" y="34" width="6" height="6"/><rect class="cls-3" x="5" y="42" width="6" height="6"/><rect class="cls-3" x="13" y="42" width="6" height="6"/><rect class="cls-3" x="21" y="42" width="6" height="6"/><rect class="cls-3" x="29" y="42" width="6" height="6"/><rect class="cls-3" x="37" y="42" width="6" height="6"/><path class="cls-2" d="M64,8H62V2H2V8H0V1A1,1,0,0,1,1,0H63a1,1,0,0,1,1,1Z"/><rect class="cls-2" x="7" y="4" width="50" height="2"/><rect class="cls-2" x="35.76" y="20" width="8.49" height="2" transform="translate(-3.13 34.44) rotate(-45)"/><rect class="cls-2" x="39" y="16.76" width="2" height="8.49" transform="translate(-3.13 34.44) rotate(-45)"/><rect class="cls-2" x="11.76" y="36" width="8.49" height="2" transform="translate(-21.48 22.15) rotate(-45)"/><rect class="cls-2" x="15" y="32.76" width="2" height="8.49" transform="translate(-21.48 22.15) rotate(-45)"/><rect class="cls-2" x="15" y="20" width="2" height="2"/><rect class="cls-2" x="7" y="20" width="2" height="2"/><rect class="cls-2" x="55" y="44" width="2" height="2"/><rect class="cls-2" x="47" y="44" width="2" height="2"/></g></g></svg>'),
       (2, 'DURATION',
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 56 64"><defs><style>.cls-1{fill:#fff;}.cls-2{fill:#4d5152;}.cls-3{fill:#67dde0;}.cls-4{fill:#f28c13;}</style></defs><g id="Layer_2" data-name="Layer 2"><g id="Layer_1-2" data-name="Layer 1"><circle class="cls-1" cx="28" cy="36" r="27"/><path class="cls-2" d="M28,64A28,28,0,1,1,56,36,28,28,0,0,1,28,64Zm0-54A26,26,0,1,0,54,36,26,26,0,0,0,28,10Z"/><path class="cls-2" d="M28,10a5,5,0,1,1,5-5A5,5,0,0,1,28,10Zm0-8a3,3,0,1,0,3,3A3,3,0,0,0,28,2Z"/><rect class="cls-2" x="27" y="5" width="2" height="7"/><circle class="cls-3" cx="28" cy="36" r="20"/><path class="cls-4" d="M28,16V36H48A20,20,0,0,0,28,16Z"/><circle class="cls-1" cx="28" cy="36" r="6"/><rect class="cls-2" x="28" y="35" width="20" height="2"/><rect class="cls-2" x="46.38" y="14.5" width="4.24" height="2" transform="translate(3.25 38.83) rotate(-45)"/><rect class="cls-2" x="6.5" y="13.38" width="2" height="4.25" transform="translate(-8.76 9.84) rotate(-45)"/><rect class="cls-2" x="49" y="11.17" width="2" height="5.66" transform="translate(4.75 39.46) rotate(-45)"/><rect class="cls-2" x="3.17" y="13" width="5.66" height="2" transform="translate(-8.14 8.34) rotate(-45)"/><circle class="cls-2" cx="28" cy="36" r="2"/><rect class="cls-2" x="1" y="35" width="3" height="2"/><rect class="cls-2" x="52" y="35" width="3" height="2"/><rect class="cls-2" x="27" y="60" width="2" height="3"/></g></g></svg>'),
       (2, 'MEET ON LOCATION',
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 64"><defs><style>.cls-1{fill:#67dde0;}.cls-2{fill:#acf0f2;}.cls-3{fill:#ffb957;}.cls-4{fill:#4d5152;}.cls-5{fill:#fff;}</style></defs><g id="Layer_2" data-name="Layer 2"><g id="Layer_1-2" data-name="Layer 1"><ellipse class="cls-1" cx="24" cy="54" rx="23" ry="10"/><ellipse class="cls-2" cx="24" cy="53" rx="9" ry="4"/><path class="cls-3" d="M43,20c0,14-19,33-19,33S5,34,5,20a19,19,0,0,1,38,0Z"/><path class="cls-4" d="M24,54.41l-.71-.7C22.5,52.92,4,34.25,4,20a20,20,0,0,1,40,0c0,14.25-18.5,32.92-19.29,33.71ZM24,2A18,18,0,0,0,6,20C6,31.92,20.59,48,24,51.56,27.41,48,42,31.91,42,20A18,18,0,0,0,24,2Z"/><circle class="cls-5" cx="24" cy="20" r="12"/><path class="cls-4" d="M24,28a8,8,0,1,1,8-8A8,8,0,0,1,24,28Zm0-14a6,6,0,1,0,6,6A6,6,0,0,0,24,14Z"/><path class="cls-4" d="M24,64C10.54,64,0,59.17,0,53c0-3.59,3.53-6.83,9.68-8.88L10.32,46C5.11,47.75,2,50.37,2,53c0,4.88,10.08,9,22,9s22-4.12,22-9c0-2.64-3.11-5.25-8.32-7l.63-1.89C44.47,46.17,48,49.41,48,53,48,59.17,37.46,64,24,64Z"/><rect class="cls-4" x="23" y="19" width="2" height="2"/></g></g></svg>'),
       (2, 'LANGUAGE',
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 60"><defs><style>.cls-1{fill:#fff;}.cls-2{fill:#67dde0;}.cls-3{fill:#4d5152;}.cls-4{fill:#ffb957;}</style></defs><g id="Layer_2" data-name="Layer 2"><g id="Layer_1-2" data-name="Layer 1"><polygon class="cls-1" points="42 25 63 25 63 1 23 1 23 14 23 31 36 31 42 25"/><path class="cls-2" d="M31,44c7,4,14,3,14,13v2H1V57c0-10,7-9,14-13Z"/><polygon class="cls-1" points="27 48 23 46 19 48 15 44 17 41 18 41 18 37 23 38 28 37 28 41 29 41 31 44 27 48"/><path class="cls-3" d="M46,60H0V57c0-8,4.45-9.7,9.15-11.49a35.44,35.44,0,0,0,5.35-2.38l1,1.74a39.47,39.47,0,0,1-5.64,2.51C5.13,49.18,2,50.37,2,57v1H44V57c0-6.63-3.13-7.82-7.86-9.62a39.47,39.47,0,0,1-5.64-2.51l1-1.74a35.44,35.44,0,0,0,5.35,2.38C41.55,47.3,46,49,46,57Z"/><path class="cls-3" d="M27.2,49.22,23,47.12l-4.2,2.1-5.09-5.09,2.82-4.22L23,41l6.47-1.08,2.82,4.22ZM23,44.88l3.8,1.9,2.91-2.91-1.18-1.78L23,43l-5.53-.92-1.18,1.78,2.91,2.91Z"/><rect class="cls-3" x="27" y="37" width="2" height="4"/><rect class="cls-3" x="17" y="37" width="2" height="4"/><rect class="cls-3" x="22" y="54" width="2" height="2"/><rect class="cls-3" x="22" y="50" width="2" height="2"/><path class="cls-1" d="M32,24v4c-1,7-4,9-4,9a13,13,0,0,1-10,0s-3-2-4-9V24l4-2,10,4Z"/><path class="cls-4" d="M14,24c0-6.08,2.92-9,9-9s9,2.92,9,9l-4,2-5-2-5-2Z"/><path class="cls-3" d="M23,39a14.21,14.21,0,0,1-5.45-1.11l-.1-.06c-.14-.09-3.39-2.33-4.44-9.69V24h2v3.93c.86,5.92,3.2,8,3.52,8.21a12.1,12.1,0,0,0,9,0c.31-.25,2.66-2.3,3.52-8.21V24h2v4.14c-1.05,7.36-4.3,9.6-4.44,9.69l-.1.06A14.21,14.21,0,0,1,23,39Z"/><path class="cls-3" d="M33,24H31c0-5.53-2.47-8-8-8s-8,2.47-8,8H13c0-6.64,3.36-10,10-10S33,17.36,33,24Z"/><polygon class="cls-3" points="28.05 27.1 18.05 23.1 14.45 24.89 13.55 23.11 17.95 20.9 27.95 24.9 31.55 23.11 32.45 24.89 28.05 27.1"/><polygon class="cls-3" points="64 26 42 26 42 24 62 24 62 2 24 2 24 12 22 12 22 0 64 0 64 26"/><path class="cls-3" d="M36,32a1,1,0,0,1-.71-.29,1,1,0,0,1,0-1.42l6-6a1,1,0,0,1,1.42,1.42l-6,6A1,1,0,0,1,36,32Z"/><rect class="cls-3" x="38" y="12" width="2" height="2"/><rect class="cls-3" x="42" y="12" width="2" height="2"/><rect class="cls-3" x="46" y="12" width="2" height="2"/></g></g></svg>'),
       (2, 'PICK-UP',
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 62 58"><defs><style>.cls-1{fill:#4d5152;}.cls-2{fill:#f28c13;}.cls-3{fill:#fff;}.cls-4{fill:#67dde0;}</style></defs><g id="Layer_2" data-name="Layer 2"><g id="Layer_1-2" data-name="Layer 1"><path class="cls-1" d="M61,58H29V56H60V4H28V54H26V3a1,1,0,0,1,1-1H61a1,1,0,0,1,1,1V57A1,1,0,0,1,61,58Z"/><rect class="cls-2" x="30" y="6" width="8" height="8"/><rect class="cls-2" x="40" y="6" width="8" height="8"/><rect class="cls-2" x="50" y="6" width="8" height="8"/><rect class="cls-2" x="40" y="16" width="8" height="8"/><rect class="cls-2" x="50" y="16" width="8" height="8"/><rect class="cls-2" x="40" y="26" width="8" height="8"/><rect class="cls-2" x="50" y="26" width="8" height="8"/><rect class="cls-2" x="40" y="36" width="8" height="8"/><rect class="cls-2" x="40" y="46" width="8" height="8"/><rect class="cls-2" x="50" y="36" width="8" height="8"/><rect class="cls-2" x="50" y="46" width="8" height="8"/><polygon class="cls-3" points="37 17 1 17 1 57 12 57 12 51 26 51 26 57 37 57 37 17"/><path class="cls-1" d="M37,58H28V56h8V18H2V56h8v2H1a1,1,0,0,1-1-1V17a1,1,0,0,1,1-1H37a1,1,0,0,1,1,1V57A1,1,0,0,1,37,58Z"/><rect class="cls-4" x="4" y="20" width="6" height="4"/><rect class="cls-4" x="12" y="20" width="6" height="4"/><rect class="cls-4" x="4" y="26" width="6" height="4"/><rect class="cls-4" x="12" y="26" width="6" height="4"/><rect class="cls-4" x="4" y="32" width="6" height="4"/><rect class="cls-4" x="12" y="32" width="6" height="4"/><rect class="cls-4" x="4" y="38" width="6" height="4"/><rect class="cls-4" x="12" y="38" width="6" height="4"/><rect class="cls-4" x="4" y="44" width="6" height="4"/><rect class="cls-4" x="12" y="44" width="6" height="4"/><rect class="cls-4" x="20" y="20" width="6" height="4"/><rect class="cls-4" x="28" y="20" width="6" height="4"/><rect class="cls-4" x="20" y="26" width="6" height="4"/><rect class="cls-4" x="28" y="26" width="6" height="4"/><rect class="cls-4" x="20" y="32" width="6" height="4"/><rect class="cls-4" x="28" y="32" width="6" height="4"/><rect class="cls-4" x="20" y="38" width="6" height="4"/><rect class="cls-4" x="28" y="38" width="6" height="4"/><rect class="cls-4" x="20" y="44" width="6" height="4"/><rect class="cls-4" x="28" y="44" width="6" height="4"/><rect class="cls-4" x="4" y="50" width="6" height="4"/><rect class="cls-4" x="28" y="50" width="6" height="4"/><rect class="cls-1" y="12" width="24" height="2"/><rect class="cls-1" y="8" width="24" height="2"/><rect class="cls-1" y="4" width="24" height="2"/><rect class="cls-1" width="24" height="2"/><path class="cls-1" d="M25,58H13a1,1,0,0,1-1-1V51a1,1,0,0,1,1-1H25a1,1,0,0,1,1,1v6A1,1,0,0,1,25,58ZM14,56H24V52H14Z"/></g></g></svg>')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS tour_features
(
    id         INTEGER PRIMARY KEY autoincrement,
    tour_id    INTEGER NOT NULL,
    feature_id INTEGER NOT NULL,
    text       TEXT DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx__tour_features__tour_id_feature_id ON tour_features (tour_id, feature_id);

CREATE TABLE IF NOT EXISTS faq_general
(
    id           INTEGER PRIMARY KEY autoincrement,
    tour_type_id SMALLINT NOT NULL,
    question     TEXT,
    answer       TEXT,
    FOREIGN KEY (tour_type_id) REFERENCES tour_types (id)
);

CREATE TABLE IF NOT EXISTS faq_tours
(
    id       INTEGER PRIMARY KEY autoincrement,
    tour_id  INTEGER,
    question TEXT,
    answer   TEXT,
    FOREIGN KEY (tour_id) REFERENCES tours (id)
);

CREATE TABLE IF NOT EXISTS highlight_types
(
    id   INTEGER PRIMARY KEY autoincrement,
    name VARCHAR NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx__highlight_types__name ON highlight_types (name);
INSERT INTO highlight_types (name)
VALUES ('highlight'),
       ('included')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS tour_highlights
(
    id                INTEGER PRIMARY KEY autoincrement,
    highlight_type_id INTEGER,
    tour_id           INTEGER,
    text              TEXT,
    FOREIGN KEY (highlight_type_id) REFERENCES highlight_types (id),
    FOREIGN KEY (tour_id) REFERENCES tours (id)
);
