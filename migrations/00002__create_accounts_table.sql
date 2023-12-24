CREATE TABLE accounts(
   account_id VARCHAR(80) PRIMARY KEY,
   user_id VARCHAR(80) NOT NULL,
   name VARCHAR(255) NOT NULL,
   balance float8 NOT NULL DEFAULT 0.0,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   deleted_at TIMESTAMPTZ NULL
);