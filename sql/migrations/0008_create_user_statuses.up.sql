CREATE TABLE user_statuses
(
    username VARCHAR(64) PRIMARY KEY NOT NULL COLLATE NOCASE,
    locked   INTEGER                 NOT NULL,

    CHECK (LENGTH(username) > 0),
    CHECK (LENGTH(username) <= 64),

    CHECK (locked IN (0, 1))
);