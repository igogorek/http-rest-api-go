package store

import (
	"fmt"
	"strings"
	"testing"
)

func TestStore(t *testing.T, databaseURL string) (*Store, func(...string)) {
	t.Helper()

	config := NewConfig()
	config.DatabaseURL = databaseURL

	st := New(config)
	if err := st.Open(); err != nil {
		t.Fatal(err)
	}
	return st, func(tables ...string) {
		if len(tables) > 0 {
			if _, err := st.db.Exec(
				fmt.Sprintf("TRUNCATE %v CASCADE", strings.Join(tables, ", ")),
			); err != nil {
				t.Fatal(err)
			}
		}
		st.Close()
	}
}
