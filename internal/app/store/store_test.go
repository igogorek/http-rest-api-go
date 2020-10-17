package store_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("GO_RESTAPI_DB_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost dbname=go_restapi_test sslmode=disable"
	}
	os.Exit(m.Run())
}
