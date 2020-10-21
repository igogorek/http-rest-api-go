package apiserver

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/igogorek/http-rest-api-go/internal/app/store"
	"github.com/igogorek/http-rest-api-go/internal/app/store/sqlstore"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// APIServer object entry point to start api server
type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  store.Store
}

// New construct object of APIServer
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start function to start APIServer
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	s.logger.Info("Starting apiserver")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
}

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.WriteString(w, "Hello"); err != nil {
			s.logger.Fatal(err)
		}
	}
}

func (s *APIServer) configureStore() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}

	s.store = sqlstore.New(db)
	return nil
}
