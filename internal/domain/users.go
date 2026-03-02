package domain

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	PublicID  string    `json:"public_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int32     `json:"version"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func NewUserFromDB(id uuid.UUID, publicID, username, email string, passwordhash []byte, activated bool, version int32) *User {
	return &User{
		ID:        id,
		PublicID:  publicID,
		Username:  username,
		Email:     email,
		Password:  password{hash: passwordhash},
		Activated: activated,
		Version:   version,
	}
}

func (u *User) PasswordHash() []byte {
	return u.Password.hash
}
