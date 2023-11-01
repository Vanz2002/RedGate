package token

import (
	"time"

	"github.com/google/uuid"
)

// struct of payload
type Payload struct {
	ID        uuid.UUID `json:"id"`
	AccountID string    `json:"account_id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(account_id string, username string, duration time.Duration) *Payload {
	token_id, err := uuid.NewRandom()
	if err != nil {
		return nil
	}

	return &Payload{
		ID:        token_id,
		AccountID: account_id,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
}

func (p *Payload) TimeValid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
