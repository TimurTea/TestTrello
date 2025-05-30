CREATE TABLE boards(
    id    SERIAL PRIMARY KEY,
    title TEXT NOT NULL
);
CREATE TABLE lists(
    id       SERIAL PRIMARY KEY,
    title    TEXT    NOT NULL,
    board_id INTEGER NOT NULL REFERENCES boards (id) ON DELETE CASCADE
);
CREATE TABLE cards(
    id       SERIAL PRIMARY KEY,
    title    TEXT    NOT NULL,
    description TEXT,
    list_id INTEGER NOT NULL REFERENCES lists (id) ON DELETE CASCADE
);