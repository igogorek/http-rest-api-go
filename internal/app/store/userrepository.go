package store

import "github.com/igogorek/http-rest-api-go/internal/app/model"

type UserRepository struct {
	store *Store
}

func (ur *UserRepository) Create(u *model.User) (*model.User, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := u.BeforeCreate(); err != nil {
		return nil, err
	}

	if err := ur.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID); err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) FindByEmail(email string) (*model.User, error) {
	user := model.User{}
	if err := ur.store.db.QueryRow(
		"SELECT * FROM users WHERE email = $1",
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.EncryptedPassword,
	); err != nil {
		return nil, err
	}
	return &user, nil
}
