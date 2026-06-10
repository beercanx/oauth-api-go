-- TODO - Remove if this is used in a production setting
INSERT INTO client_configurations (client_id, client_type, redirect_uris, allowed_scopes, allowed_actions, allowed_grant_types)
VALUES ('aardvark', 'confidential', '[]', '["basic", "read", "write"]', '["introspect"]', '["password"]'),
       ('cicada', 'public', '["https://cicada.baconi.co.uk/callback"]', '["basic"]', '["authorise", "pkce"]', '["authorization_code"]'),
       ('dodo', 'confidential', '["https://dodo.baconi.co.uk/callback"]', '["basic"]', '["authorise", "pkce"]', '["authorization_code"]');
