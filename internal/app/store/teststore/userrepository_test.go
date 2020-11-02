package teststore_test

import (
	"github.com/igogorek/http-rest-api-go/internal/app/model"
	"github.com/igogorek/http-rest-api-go/internal/app/store"
	"github.com/igogorek/http-rest-api-go/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	st := teststore.New()

	user := model.TestUser()
	err := st.User().Create(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.EncryptedPassword)
}

func TestUserRepository_Find(t *testing.T) {
	st := teststore.New()

	user, err := st.User().Find(1)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, user)

	testUser := model.TestUser()
	if err = st.User().Create(testUser); err != nil {
		t.Fatal(err)
	}

	user, err = st.User().Find(testUser.ID)

	assert.NoError(t, err)
	assert.Equal(t,
		model.User{
			ID:                testUser.ID,
			Email:             testUser.Email,
			Password:          testUser.Password,
			EncryptedPassword: testUser.EncryptedPassword,
		},
		*user,
	)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	st := teststore.New()

	testUser := model.TestUser()

	user, err := st.User().FindByEmail(testUser.Email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, user)

	if err = st.User().Create(testUser); err != nil {
		t.Fatal(err)
	}

	user, err = st.User().FindByEmail(testUser.Email)

	assert.NoError(t, err)
	assert.Equal(t,
		model.User{
			ID:                testUser.ID,
			Email:             testUser.Email,
			Password:          testUser.Password,
			EncryptedPassword: testUser.EncryptedPassword,
		},
		*user,
	)
}
