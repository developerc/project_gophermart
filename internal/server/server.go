package server

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	dbstorage "github.com/developerc/project_gophermart/internal/db_storage"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type svc interface {
	Register(buf bytes.Buffer) (*http.Cookie, error)
	UserLogin(buf bytes.Buffer) (*http.Cookie, error)
}

type Server struct {
	service svc
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var err error
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cookie, err := s.service.Register(buf)
	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) && pgErr.ConstraintName == "must_be_different_usr":
			http.Error(w, err.Error(), http.StatusConflict)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	fmt.Println("cookie:", cookie)
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

func (s *Server) UserLogin(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var err error
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cookie, err := s.service.UserLogin(buf)
	if err != nil {
		if _, ok := err.(*dbstorage.ErrorLgnPsw); ok {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(cookie)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func NewServer(service svc) (*Server, error) {
	srv := new(Server)
	srv.service = service
	return srv, nil
}

func (s *Server) SetupRoutes() http.Handler {
	r := chi.NewRouter()
	r.Post("/api/user/register", s.Register)
	r.Post("/api/user/login", s.UserLogin)
	return r
}
