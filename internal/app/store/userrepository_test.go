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

	user, err := st.User().Create(model.TestUser())

	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	st, teardown := store.TestStore(t, databaseURL)
	defer teardown("users")

	test_user := model.TestUser()

	user, err := st.User().FindByEmail(test_user.Email)
	assert.Error(t, err)
	assert.Nil(t, user)

	created, err := st.User().Create(test_user)

	if err != nil {
		t.Fatal(err)
	}

	user, err = st.User().FindByEmail(test_user.Email)

	assert.NoError(t, err)
	assert.Equal(t,
		model.User{
			ID:                created.ID,
			Email:             created.Email,
			EncryptedPassword: created.EncryptedPassword,
		},
		*user,
	)
}
