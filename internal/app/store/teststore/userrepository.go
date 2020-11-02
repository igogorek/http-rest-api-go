package teststore

import (
	"github.com/igogorek/http-rest-api-go/internal/app/model"
	"github.com/igogorek/http-rest-api-go/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[int]*model.User
}

func (ur *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	u.ID = len(ur.users) + 1
	ur.users[u.ID] = u

	return nil
}

func (ur *UserRepository) Find(id int) (*model.User, error) {
	user, ok := ur.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return user, nil
}

func (ur *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, user := range ur.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, store.ErrRecordNotFound

}
