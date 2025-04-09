-- name: CreateAccount :exec
INSERT INTO accounts (
    user_id, first_name, last_name, username, phone, about, birthday, personal_channel_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) ON CONFLICT (user_id) DO NOTHING;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE user_id = $1
LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY created_at
LIMIT $1
OFFSET $2;

-- name: AccountsQuantity :one
SELECT COUNT(*) FROM accounts;
