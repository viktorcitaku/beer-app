CREATE SCHEMA beer;

CREATE TABLE IF NOT EXISTS beer.user_profile
(
    email       VARCHAR(255) PRIMARY KEY,
    last_update TIMESTAMP WITHOUT TIME ZONE
);

INSERT INTO beer.user_profile(email, last_update) VALUES ('palodo3046@fom8.com', current_timestamp);
INSERT INTO beer.user_profile(email, last_update) VALUES ('nolessam@pianoxltd.com', current_timestamp);

COMMIT;