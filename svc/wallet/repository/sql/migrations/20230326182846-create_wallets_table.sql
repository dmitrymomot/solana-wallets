
-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION wallets_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TABLE IF NOT EXISTS wallets (
    user_id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL DEFAULT 'default',
    public_key VARCHAR NOT NULL,
    mnemonic text NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NULL
);
CREATE UNIQUE INDEX wallets_public_key ON wallets (public_key);
CREATE TRIGGER update_wallets_modtime BEFORE
UPDATE ON wallets FOR EACH ROW EXECUTE PROCEDURE wallets_update_updated_at_column();
-- +migrate StatementEnd

-- +migrate Down
DROP TRIGGER IF EXISTS update_wallets_modtime ON wallets;
DROP TABLE IF EXISTS wallets;
DROP FUNCTION IF EXISTS wallets_update_updated_at_column();