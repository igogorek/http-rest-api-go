package apiserver

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/igogorek/http-rest-api-go/internal/app/model"
	"github.com/igogorek/http-rest-api-go/internal/app/store"
	"github.com/sirupsen/logrus"
	"net/http"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	srv := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	srv.configureRouter()

	return srv
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/users", s.handleUserCreate()).Methods(http.MethodPost)
}

func (s *server) configureLogger(logLevel string) error {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) handleError(
	w http.ResponseWriter,
	r *http.Request,
	statusCode int,
	err error,
) {
	s.respond(w, r, statusCode, err)
}

func (s *server) respond(
	w http.ResponseWriter,
	r *http.Request,
	statusCode int,
	data interface{},
) {
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			s.logger.Error(err)
		}
	}
}

func (s *server) handleUserCreate() http.HandlerFunc {
	type reqData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := reqData{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			s.handleError(w, r, http.StatusBadRequest, err)
			return
		}

		user := &model.User{
			Email:    data.Email,
			Password: data.Password,
		}

		if err := s.store.User().Create(user); err != nil {
			s.handleError(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		user.Sanitize()
		s.respond(w, r, http.StatusCreated, user)
	}
}
