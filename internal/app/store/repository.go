package store

import "github.com/igogorek/http-rest-api-go/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
}
