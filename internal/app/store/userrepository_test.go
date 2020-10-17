package store_test

import (
	"github.com/igogorek/http-rest-api-go/internal/app/model"
	"github.com/igogorek/http-rest-api-go/internal/app/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	st, teardown := store.TestStore(t, databaseURL)
	defer teardown("users")

	user, err := st.User().Create(&model.User{
		Email: "example@example.org",
	})

	assert.NoError(t, err)
	assert.NotNil(t, user)
}
