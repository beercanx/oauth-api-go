CREATE TABLE access_tokens
(
    id         VARCHAR(36) PRIMARY KEY NOT NULL,
    username   VARCHAR(64)             NOT NULL,
    client_id  VARCHAR(64)             NOT NULL,
    scopes     JSONB                   NOT NULL,
    issued_at  TIMESTAMP               NOT NULL,
    expires_at TIMESTAMP               NOT NULL,
    not_before TIMESTAMP               NOT NULL,

    CHECK (LENGTH(id) = 36),
    CHECK (id != '00000000-0000-0000-0000-000000000000'),

    CHECK (LENGTH(username) > 0),
    CHECK (LENGTH(username) <= 64),

    CHECK (LENGTH(client_id) > 0),
    CHECK (LENGTH(client_id) <= 64),

    CHECK (json_valid(scopes)),
    CHECK (json_type(scopes) = 'array'),

    CHECK (issued_at != '0001-01-01T00:00:00Z'),
    CHECK (expires_at != '0001-01-01T00:00:00Z'),
    CHECK (not_before != '0001-01-01T00:00:00Z')
);
