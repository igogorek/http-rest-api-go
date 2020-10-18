package model_test

import (
	"github.com/igogorek/http-rest-api-go/internal/app/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_BeforeCreate(t *testing.T) {
	u := model.TestUser()
	err := u.BeforeCreate()

	assert.NoError(t, err)
	assert.NotEmpty(t, u.EncryptedPassword)
}
