package teststore

import (
	"github.com/igogorek/http-rest-api-go/internal/app/model"
	"github.com/igogorek/http-rest-api-go/internal/app/store"
)

type Store struct {
	userRepository *UserRepository
}

func New() store.Store {
	return &Store{}
}

func (st *Store) User() store.UserRepository {
	if st.userRepository == nil {
		st.userRepository = &UserRepository{
			store: st,
			users: make(map[string]*model.User),
		}
	}

	return st.userRepository
}
