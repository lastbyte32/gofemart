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

CREATE TABLE  orders
(
	number      TEXT  NOT NULL PRIMARY KEY,
	user_id     UUID NOT NULL,
	status      TEXT  NOT NULL,
	accrual 	FLOAT NOT NULL DEFAULT 0,
	uploaded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE withdrawals
(
		id SERIAL,
		user_id UUID NOT NULL,
		order_number text NOT NULL,
		sum float NOT NULL,
		processed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT accruals_withdrawn_pkey PRIMARY KEY (id)
);

COMMIT;
