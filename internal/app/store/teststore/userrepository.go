package teststore

import (
	"github.com/igogorek/http-rest-api-go/internal/app/model"
	"github.com/igogorek/http-rest-api-go/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[string]*model.User
}

func (ur *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	ur.users[u.Email] = u
	u.ID = len(ur.users)

	return nil
}

func (ur *UserRepository) FindByEmail(email string) (*model.User, error) {
	user, ok := ur.users[email]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return user, nil
}
