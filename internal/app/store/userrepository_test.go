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

func TestUserRepository_FindByEmail(t *testing.T) {
	st, teardown := store.TestStore(t, databaseURL)
	defer teardown("users")

	email := "example@example.org"

	user, err := st.User().FindByEmail(email)
	assert.Error(t, err)
	assert.Nil(t, user)

	created, err := st.User().Create(&model.User{
		Email: email,
	})

	if err != nil {
		t.Fatal(err)
	}

	user, err = st.User().FindByEmail(email)

	assert.NoError(t, err)
	assert.Equal(t, created, user)
}
