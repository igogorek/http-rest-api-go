package store

import "github.com/igogorek/http-rest-api-go/internal/app/model"

type UserRepository struct {
	store *Store
}

func (ur *UserRepository) Create(u *model.User) (*model.User, error) {
	if err := ur.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID); err != nil {
		return nil, err
	}
	return u, nil
}
