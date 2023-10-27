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
  email, password_hash
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteAccount :one
DELETE FROM account
WHERE account_id = $1
RETURNING *;