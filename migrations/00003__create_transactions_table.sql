CREATE TABLE transactions(
    transaction_id VARCHAR(80) NOT NULL PRIMARY KEY,
    account_id VARCHAR(80) NOT NULL,
    user_id VARCHAR(80) NOT NULL,
    amount float8 NOT NULL,
    balance float8 NOT NULL,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(80) NOT NULL,
    metadata JSON,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);