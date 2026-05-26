CREATE TABLE access_tokens
(
    id         VARCHAR(36) PRIMARY KEY NOT NULL CHECK (LENGTH(id) = 36),
    username   VARCHAR(64)             NOT NULL CHECK (LENGTH(username) > 0 AND LENGTH(username) <= 64),
    client_id  VARCHAR(64)             NOT NULL CHECK (LENGTH(client_id) > 0 AND LENGTH(client_id) <= 64),
    scopes     TEXT                    NOT NULL,
    issued_at  timestamp               NOT NULL,
    expires_at timestamp               NOT NULL,
    not_before timestamp               NOT NULL
);
