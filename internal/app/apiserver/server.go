package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/igogorek/http-rest-api-go/internal/app/model"
	"github.com/igogorek/http-rest-api-go/internal/app/store"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	sessionName      = "apiserver-session"
	sessionKeyUserId = "user_id"
)
const (
	ctxKeyUser ctxKey = iota
	ctxRequestId
)

var (
	errInternalError          = errors.New("internal error")
	errInvalidEmailOrPassword = errors.New("invalid email or password")
	errUnauthenticatedUser    = errors.New("unauthenticated user")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	srv := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	srv.configureRouter()

	return srv
}

func (s *server) configureRouter() {
	s.router.Use(s.addRequestId)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/users", s.handleUserCreate()).Methods(http.MethodPost)
	s.router.HandleFunc("/sessions", s.handleSessionCreate()).Methods(http.MethodPost)

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods(http.MethodGet)
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
	s.respond(w, r, statusCode, map[string]string{"error": err.Error()})
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

func (s *server) addRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuidStr := uuid.New().String()
		w.Header().Set("X-Request-Id", uuidStr)
		uidContext := context.WithValue(r.Context(), ctxRequestId, uuidStr)
		next.ServeHTTP(w, r.WithContext(uidContext))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxRequestId),
		})

		logger.Infof("Starting request %s %s", r.Method, r.RequestURI)

		start := time.Now()
		codeRW := &ResponseWriter{w, http.StatusOK}
		next.ServeHTTP(codeRW, r)

		logger.Infof(
			"request ended with (%d %s) in %v",
			codeRW.code,
			http.StatusText(codeRW.code),
			time.Now().Sub(start),
		)
	})
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

func (s *server) handleSessionCreate() http.HandlerFunc {
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

		user, err := s.store.User().FindByEmail(data.Email)
		if err != nil || !user.ValidPassword(data.Password) {
			s.handleError(w, r, http.StatusUnauthorized, errInvalidEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.handleError(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values[sessionKeyUserId] = user.ID
		if err := session.Save(r, w); err != nil {
			s.handleError(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}

}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.handleError(w, r, http.StatusInternalServerError, err)
			return
		}

		userId, ok := session.Values[sessionKeyUserId]
		if !ok {
			s.handleError(w, r, http.StatusUnauthorized, errUnauthenticatedUser)
			return
		}

		user, err := s.store.User().Find(userId.(int))
		if err != nil {
			s.handleError(w, r, http.StatusUnauthorized, errUnauthenticatedUser)
			return
		}

		userContext := context.WithValue(r.Context(), ctxKeyUser, user)
		next.ServeHTTP(w, r.WithContext(userContext))
	})
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ctxKeyUser).(*model.User)
		if !ok {
			s.handleError(w, r, http.StatusInternalServerError, errInternalError)
			return
		}
		s.respond(w, r, http.StatusOK, user)
	}
}
