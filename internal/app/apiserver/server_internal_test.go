package apiserver

import (
	"bytes"
	"encoding/json"
	"github.com/igogorek/http-rest-api-go/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleUserCreate(t *testing.T) {
	srv := newServer(teststore.New())

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "user@example.org",
				"password": "password",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "cant be decoded as email password",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"email":    "invalid@email",
				"password": "password",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    "valid@email.org",
				"password": "short",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "missing email",
			payload: map[string]string{
				"password": "password",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "missing password",
			payload: map[string]string{
				"email": "valid@email.org",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := bytes.NewBuffer([]byte{})
			if err := json.NewEncoder(b).Encode(tc.payload); err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/users", b)
			if err != nil {
				t.Fatal(err)
			}

			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

}
