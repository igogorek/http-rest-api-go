package sqlstore_test

import (
	"github.com/igogorek/http-rest-api-go/internal/app/model"
	"github.com/igogorek/http-rest-api-go/internal/app/store"
	"github.com/igogorek/http-rest-api-go/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	st := sqlstore.New(db)

	user := model.TestUser()
	err := st.User().Create(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.EncryptedPassword)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	st := sqlstore.New(db)

	test_user := model.TestUser()

	user, err := st.User().FindByEmail(test_user.Email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, user)

	if err = st.User().Create(test_user); err != nil {
		t.Fatal(err)
	}

	user, err = st.User().FindByEmail(test_user.Email)

	assert.NoError(t, err)
	assert.Equal(t,
		model.User{
			ID:                test_user.ID,
			Email:             test_user.Email,
			EncryptedPassword: test_user.EncryptedPassword,
		},
		*user,
	)
}
