package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID                int
	Email             string
	Password          string
	EncryptedPassword string
}

func (u *User) BeforeCreate() error {
	if u.Password == "" {
		return nil
	}

	enc, err := encryptString(u.Password)
	if err != nil {
		return err
	}

	u.EncryptedPassword = enc
	return nil
}

func encryptString(s string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
