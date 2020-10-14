package apiserver

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIServer_HandleHello(t *testing.T) {
	s := New(NewConfig())

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/hello", nil)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	s.handleHello().ServeHTTP(rec, req)
	assert.Equal(t, "Hello", rec.Body.String())
}
