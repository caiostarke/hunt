package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	AvailableUsers []*User
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Level          int       `json:"level"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	ProfilePicture string    `json:"profile_picture"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	Password       Password  `json:"-"`
}

// You can instead use this structure declare an anonymous struct inside Handler
type UserJSON struct {
	ID             uuid.UUID `json:"id"`
	ProfilePicture string    `json:"profile_picture"`
	Name           string    `json:"name" binding:"required"`
	Email          string    `json:"email" binding:"required"`
	Password       string    `json:"password" binding:"required"`
}

type Password struct {
	Plaintext string
	Hash      []byte
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.Plaintext = plaintextPassword
	p.Hash = hash
	return nil
}

func (p *Password) CheckPasswordHash() bool {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(p.Plaintext))
	return err == nil
}

type UserID struct {
	ID uuid.UUID `json:"id"`
}

func NewUser() *User {
	return &User{
		ID: uuid.New(),
	}
}
