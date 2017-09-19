package data

import (
	"database/sql"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Token
type Token struct {
	ID      string
	UserID  string
	Created time.Time
}

// TokenService is our interface for what a concrete service should implement
type TokenStore interface {
	New(uid string) (Token, error)
	ToAuth(tk Token) (string, error)
}

// TokenService is our concrete implementation of the TokenService
type TokenService struct {
	db  *sql.DB
	key []byte
}

// InitTokenService returns a token service using the given key
func InitTokenService(key string, db *sql.DB) *TokenService {
	return &TokenService{key: []byte(key), db: db}
}

// New creates a new Token in the database. It returns the created token object, signed jwt and/or an error
func (ts *TokenService) New(uid string) (Token, error) {
	id, err := uuid.NewRandom()

	if err != nil {
		return Token{}, err
	}

	t := Token{
		ID:      id.String(),
		UserID:  uid,
		Created: time.Now(),
	}

	_, err = ts.db.Exec(
		"Insert into tokens (token_id, user_id, date_created) VALUES ($1, $2, $3)",
		t.ID,
		t.UserID,
		t.Created,
	)

	if err != nil {
		return Token{}, err
	}

	return t, nil
}

// New returns a new key with the given token id
func (ts *TokenService) ToAuth(tk Token) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token_id": tk.ID,
		"user_id":  tk.UserID,
	})

	return token.SignedString(ts.key)
}

// claims returns our standard jwt claims
func claims() *jwt.StandardClaims {
	return &jwt.StandardClaims{
		ExpiresAt: 1916002709, // Expire in 2030, this should probably be changed though :)
		Issuer:    "delay",
	}
}
