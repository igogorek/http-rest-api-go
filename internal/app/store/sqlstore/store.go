package sqlstore

import (
	"database/sql"
	"github.com/igogorek/http-rest-api-go/internal/app/store"
	_ "github.com/lib/pq"
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

func New(db *sql.DB) store.Store {
	return &Store{
		db: db,
	}
}

func (st *Store) User() store.UserRepository {
	if st.userRepository == nil {
		st.userRepository = &UserRepository{
			store: st,
		}
	}

	return st.userRepository
}
