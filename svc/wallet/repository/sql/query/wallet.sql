-- name: CreateWallet :one
INSERT INTO wallets (user_id, name, public_key, mnemonic) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetWallet :one
SELECT * FROM wallets WHERE user_id = $1;

-- name: GetWalletByPublicKey :one
SELECT * FROM wallets WHERE public_key = $1;

-- name: UpdateWallet :one
UPDATE wallets SET name = $2, mnemonic = $3 WHERE user_id = $1 RETURNING *;

-- name: DeleteWallet :exec
DELETE FROM wallets WHERE user_id = $1;