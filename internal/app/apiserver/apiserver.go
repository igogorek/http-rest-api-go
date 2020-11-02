package apiserver

import (
	"database/sql"
	"github.com/gorilla/sessions"
	"github.com/igogorek/http-rest-api-go/internal/app/store/sqlstore"
	"net/http"
	"os"
)

// Start function to start apiserver
func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(os.Getenv("APISERVER_SESSION_KEY")))

	srv := newServer(store, sessionStore)

	if err := srv.configureLogger(config.LogLevel); err != nil {
		return err
	}
	srv.logger.Infof("Starting apiserver on %v", config.BindAddr)
	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
