BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    -- id UUID DEFAULT uuid_generate_v4()  PRIMARY KEY,
    id        UUID NOT NULL PRIMARY KEY,
    login     TEXT UNIQUE NOT NULL,
    password  TEXT  NOT NULL,
    create_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT login_unique UNIQUE (login)

);
COMMIT;
