CREATE TABLE user_credentials
(
    username    VARCHAR(64) PRIMARY KEY NOT NULL REFERENCES user_statuses (username) ON DELETE CASCADE,
    hash        VARCHAR(128)            NOT NULL,
    created_at  TIMESTAMP               NOT NULL DEFAULT current_timestamp,
    modified_at TIMESTAMP               NOT NULL DEFAULT current_timestamp,

    CHECK (LENGTH(hash) >= 50 AND LENGTH(hash) <= 128),

    CHECK (created_at != '0001-01-01T00:00:00Z'),
    CHECK (modified_at != '0001-01-01T00:00:00Z')
);

CREATE TRIGGER user_credentials_modified_at
    AFTER UPDATE ON user_credentials
    FOR EACH ROW
BEGIN
    UPDATE user_credentials
    SET modified_at = current_timestamp
    WHERE username = OLD.username;
END;
