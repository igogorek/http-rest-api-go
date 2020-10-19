package model_test

import (
	"github.com/igogorek/http-rest-api-go/internal/app/model"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		prepareUser func() *model.User
		isValid     bool
	}{
		{
			name: "valid",
			prepareUser: func() *model.User {
				return model.TestUser()
			},
			isValid: true,
		},
		{
			name: "emptyEmail",
			prepareUser: func() *model.User {
				user := model.TestUser()
				user.Email = ""

				return user
			},
			isValid: false,
		},
		{
			name: "invalidEmail",
			prepareUser: func() *model.User {
				user := model.TestUser()
				user.Email = "invalid@email"

				return user
			},
			isValid: false,
		},
		{
			name: "emptyPassword",
			prepareUser: func() *model.User {
				user := model.TestUser()
				user.Password = ""

				return user
			},
			isValid: false,
		},
		{
			name: "tooShortPassword",
			prepareUser: func() *model.User {
				user := model.TestUser()
				user.Password = "12345"

				return user
			},
			isValid: false,
		},
		{
			name: "shortPassword",
			prepareUser: func() *model.User {
				user := model.TestUser()
				user.Password = "123456"

				return user
			},
			isValid: true,
		},
		{
			name: "tooLongPassword",
			prepareUser: func() *model.User {
				user := model.TestUser()
				var b strings.Builder
				for i := 0; i < 10; i++ {
					b.Write([]byte("0123456789"))
				}
				user.Password = b.String() + "1"
				return user
			},
			isValid: false,
		},
		{
			name: "longPassword",
			prepareUser: func() *model.User {
				user := model.TestUser()
				var b strings.Builder
				for i := 0; i < 10; i++ {
					b.Write([]byte("0123456789"))
				}
				user.Password = b.String()
				return user
			},
			isValid: true,
		},
		{
			name: "encryptedPassword",
			prepareUser: func() *model.User {
				user := model.TestUser()
				user.Password = ""
				user.EncryptedPassword = "asdfadfadsf"
				return user
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.prepareUser().Validate())
			} else {
				assert.Error(t, tc.prepareUser().Validate())
			}
		})
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	u := model.TestUser()
	err := u.BeforeCreate()

	assert.NoError(t, err)
	assert.NotEmpty(t, u.EncryptedPassword)
}
