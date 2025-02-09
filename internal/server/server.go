package server

import (
	"bytes"
	"errors"
	"log"

	//"fmt"
	"net/http"

	dbstorage "github.com/developerc/project_gophermart/internal/db_storage"
	"github.com/developerc/project_gophermart/internal/general"
	"github.com/developerc/project_gophermart/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type svc interface {
	Register(buf bytes.Buffer) (*http.Cookie, error)
	UserLogin(buf bytes.Buffer) (*http.Cookie, error)
	GetUserFromCookie(cookieValue string) (string, error)
	PostUserOrders(usr string, buf bytes.Buffer) error
	GetUserOrders(usr string) ([]byte, error)
	GetUserBalance(usr string) ([]byte, error)
	PostBalanceWithdraw(usr string, buf bytes.Buffer) error
	GetUserWithdrawals(usr string) ([]byte, error)
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
	//fmt.Println("cookie:", cookie)
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
	//fmt.Println(cookie)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) PostUserOrders(w http.ResponseWriter, r *http.Request) {
	var usr string
	var buf bytes.Buffer
	//проверим ессть ли куки, узнаем юзера
	cookie, err := r.Cookie("user")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	usr, err = s.service.GetUserFromCookie(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Println("usr:", usr)
	//заберем боди
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.service.PostUserOrders(usr, buf)
	//много вариантов err
	if err != nil {
		if _, ok := err.(*general.ErrorNumOrder); ok {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if _, ok := err.(*general.ErrorExistsOrderSame); ok {
			http.Error(w, err.Error(), http.StatusOK)
			return
		}
		if _, ok := err.(*general.ErrorExistsOrderOther); ok {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("from GetUserOrders")
	var usr string
	var jsonBytes []byte
	//проверим ессть ли куки, узнаем юзера
	cookie, err := r.Cookie("user")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	usr, err = s.service.GetUserFromCookie(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("from GetUserOrders usr: ", usr)
	jsonBytes, err = s.service.GetUserOrders(usr)
	if err != nil {
		if _, ok := err.(*general.ErrorNoContent); ok {
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonBytes); err != nil {
		return
	}
}

func (s *Server) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("from GetUserBalance")
	var usr string
	var jsonBytes []byte
	cookie, err := r.Cookie("user")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	usr, err = s.service.GetUserFromCookie(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonBytes, err = s.service.GetUserBalance(usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonBytes); err != nil {
		return
	}
}

func (s *Server) PostBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("from PostBalanceWithdraw")
	var usr string
	var buf bytes.Buffer
	cookie, err := r.Cookie("user")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	usr, err = s.service.GetUserFromCookie(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Println("usr:", usr)
	//заберем боди
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.service.PostBalanceWithdraw(usr, buf)
	if err != nil {
		/*if _, ok := err.(*general.ErrorNumOrder); ok {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}*/
		if _, ok := err.(*general.ErrorLoyaltyPoints); ok {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return
		}
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) && pgErr.ConstraintName == "must_be_different_order":
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func (s *Server) GetUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("from GetUserWithdrawals")
	var usr string
	cookie, err := r.Cookie("user")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	usr, err = s.service.GetUserFromCookie(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Println("usr:", usr)
	jsonBytes, err := s.service.GetUserWithdrawals(usr)
	if err != nil {
		if _, ok := err.(*general.ErrorNoContent); ok {
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonBytes); err != nil {
		return
	}
}

func NewServer(service svc) (*Server, error) {
	srv := new(Server)
	srv.service = service
	return srv, nil
}

func (s *Server) SetupRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.GzipHandle)
	r.Post("/api/user/register", s.Register)
	r.Post("/api/user/login", s.UserLogin)
	r.Post("/api/user/orders", s.PostUserOrders)
	r.Get("/api/user/orders", s.GetUserOrders)
	r.Get("/api/user/balance", s.GetUserBalance)
	r.Post("/api/user/balance/withdraw", s.PostBalanceWithdraw)
	r.Get("/api/user/withdrawals", s.GetUserWithdrawals)
	return r
}
