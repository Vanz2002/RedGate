-- name: GetAccount :one
SELECT * FROM account
WHERE account_id = $1 LIMIT 1;

-- name: GetAccountbyEmail :one
SELECT * FROM account
WHERE email = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM account
ORDER BY created_at ASC;

-- name: CreateAccount :one
INSERT INTO account (
  username, email, password_hash
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteAccount :one
DELETE FROM account
WHERE account_id = $1
RETURNING *;

-- name: AccountSubscribe :one
UPDATE account SET is_subscribe = true WHERE account_id = $1
RETURNING is_subscribe;