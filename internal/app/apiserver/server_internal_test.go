package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/igogorek/http-rest-api-go/internal/app/model"
	"github.com/igogorek/http-rest-api-go/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleUserCreate(t *testing.T) {
	srv := newServer(teststore.New(), sessions.NewCookieStore([]byte("simple")))

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

func TestServer_HandleSessionCreate(t *testing.T) {
	store := teststore.New()
	user := model.TestUser()

	if err := store.User().Create(user); err != nil {
		t.Fatal(err)
	}

	srv := newServer(store, sessions.NewCookieStore([]byte("simple")))

	testCases := []struct {
		name     string
		payload  interface{}
		httpCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    user.Email,
				"password": user.Password,
			},
			httpCode: http.StatusOK,
		},
		{
			name:     "invalid payload",
			payload:  "invalid",
			httpCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"email":    "invalid",
				"password": user.Password,
			},
			httpCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    user.Email,
				"password": "invalid",
			},
			httpCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := bytes.NewBuffer([]byte{})
			if err := json.NewEncoder(b).Encode(tc.payload); err != nil {
				t.Fatal(err)
			}

			request := httptest.NewRequest(http.MethodPost, "/sessions", b)
			srv.ServeHTTP(rec, request)

			assert.Equal(t, tc.httpCode, rec.Code)
		})
	}
}

func TestServer_AuthenticateUser(t *testing.T) {
	store := teststore.New()
	user := model.TestUser()
	if err := store.User().Create(user); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name        string
		cookieValue map[interface{}]interface{}
		httpCode    int
	}{
		{
			name: "authenticated",
			cookieValue: map[interface{}]interface{}{
				sessionKeyUserId: user.ID,
			},
			httpCode: http.StatusOK,
		},
		{
			name:        "not authenticated",
			cookieValue: nil,
			httpCode:    http.StatusUnauthorized,
		},
	}

	secretKey := []byte("simple")
	srv := newServer(store, sessions.NewCookieStore(secretKey))
	sc := securecookie.New(secretKey, nil)
	mockHandler := srv.authenticateUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			cookieStr, err := sc.Encode(sessionName, tc.cookieValue)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			mockHandler.ServeHTTP(rec, req)

			assert.Equal(t, tc.httpCode, rec.Code)
		})
	}
}
