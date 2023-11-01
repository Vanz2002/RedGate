// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: query.account.sql

package sqlc

import (
	"context"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO account (
  username, email, password_hash
) VALUES (
  $1, $2, $3
)
RETURNING account_id, username, email, password_hash, created_at, is_subscribe
`

type CreateAccountParams struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount, arg.Username, arg.Email, arg.PasswordHash)
	var i Account
	err := row.Scan(
		&i.AccountID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.IsSubscribe,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :one
DELETE FROM account
WHERE account_id = $1
RETURNING account_id, username, email, password_hash, created_at, is_subscribe
`

func (q *Queries) DeleteAccount(ctx context.Context, accountID string) (Account, error) {
	row := q.db.QueryRowContext(ctx, deleteAccount, accountID)
	var i Account
	err := row.Scan(
		&i.AccountID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.IsSubscribe,
	)
	return i, err
}

const getAccount = `-- name: GetAccount :one
SELECT account_id, username, email, password_hash, created_at, is_subscribe FROM account
WHERE account_id = $1 LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, accountID string) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, accountID)
	var i Account
	err := row.Scan(
		&i.AccountID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.IsSubscribe,
	)
	return i, err
}

const getAccountbyEmail = `-- name: GetAccountbyEmail :one
SELECT account_id, username, email, password_hash, created_at, is_subscribe FROM account
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetAccountbyEmail(ctx context.Context, email string) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountbyEmail, email)
	var i Account
	err := row.Scan(
		&i.AccountID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.IsSubscribe,
	)
	return i, err
}

const listAccounts = `-- name: ListAccounts :many
SELECT account_id, username, email, password_hash, created_at, is_subscribe FROM account
ORDER BY created_at ASC
`

func (q *Queries) ListAccounts(ctx context.Context) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, listAccounts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.AccountID,
			&i.Username,
			&i.Email,
			&i.PasswordHash,
			&i.CreatedAt,
			&i.IsSubscribe,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
