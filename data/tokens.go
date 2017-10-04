package data

import (
	"database/sql"

	"time"

	"errors"

	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var ErrInvalidTokenSigningMethod = errors.New("tokens: Invalid signing method used to sign jwt")
var ErrInvalidTokenClaims = errors.New("tokens: Token jwt claims invalid")
var ErrCannotParseToken = errors.New("tokens: Cannot parse token")
var ErrTokenExpired = errors.New("tokens: Token expired")
var ErrTokenInvalid = errors.New("tokens: Given token not present in db")

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
	FromAuth(tks string) (Token, error)
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
		return Token{}, fmt.Errorf("token - New: failed to generate UUID: %v", err)
	}

	t := Token{
		ID:      id.String(),
		UserID:  uid,
		Created: time.Now().Round(time.Second),
	}

	_, err = ts.db.Exec(
		"Insert into tokens (token_id, user_id, date_created) VALUES ($1, $2, $3)",
		t.ID,
		t.UserID,
		t.Created,
	)

	if err != nil {
		return Token{}, fmt.Errorf("token - New: failed to insert token into db: %v", err)
	}

	return t, nil
}

// New returns a new key with the given token id
func (ts *TokenService) ToAuth(tk Token) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token_id": tk.ID,
	})

	return token.SignedString(ts.key)
}

// From auth takes a string jwt and returns the token struct with values from the
// database. Returns an error if it cant parse the token or the token is not present
// in the database.
func (ts *TokenService) FromAuth(tks string) (Token, error) {
	jtk, err := jwt.Parse(tks, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, ErrInvalidTokenSigningMethod
		}

		return ts.key, nil
	})

	if err != nil {
		return Token{}, ErrCannotParseToken
	}

	claims, ok := jtk.Claims.(jwt.MapClaims)

	if !ok {
		return Token{}, ErrInvalidTokenClaims
	}

	if !jtk.Valid {
		return Token{}, ErrTokenExpired
	}

	tid, ok := claims["token_id"].(string)

	if !ok {
		return Token{}, ErrCannotParseToken
	}

	dbtk := Token{}

	row := ts.db.QueryRow("SELECT token_id, user_id, date_created FROM tokens WHERE token_id = $1", tid)

	err = row.Scan(&dbtk.ID, &dbtk.UserID, &dbtk.Created)

	if err != nil {
		if err == sql.ErrNoRows {
			return Token{}, ErrTokenInvalid
		}

		return Token{}, fmt.Errorf("tokens - FromAuth: Failed to get token from db: %v", err)
	}

	return Token{
		UserID:  dbtk.UserID,
		ID:      dbtk.ID,
		Created: dbtk.Created.In(time.Local),
	}, nil
}
