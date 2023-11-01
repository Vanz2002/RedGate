package token

import "time"

// this is an functional intercafe
type Maker interface {
	GenerateToken(account_id string, username string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
