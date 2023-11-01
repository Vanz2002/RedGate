package token

import (
	"errors"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	pV2 *paseto.V2
	key []byte
}

func NewPasetoMaker(key string) (Maker, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, errors.New("invalid secret key length")
	}
	p := paseto.NewV2()
	return &PasetoMaker{p, []byte(key)}, nil
}

// VerifyToken(token string) (*Payload, error)

func (p *PasetoMaker) GenerateToken(account_id string, username string, duration time.Duration) (string, *Payload, error) {
	payload := NewPayload(account_id, username, duration)

	token, err := p.pV2.Encrypt(p.key, payload, nil)
	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

func (p *PasetoMaker) VerifyToken(tokenString string) (*Payload, error) {
	var payload Payload

	err := p.pV2.Decrypt(tokenString, p.key, &payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.TimeValid()
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
