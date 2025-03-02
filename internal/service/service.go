package service

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/developerc/project_gophermart/internal/config"
	dbstorage "github.com/developerc/project_gophermart/internal/db_storage"
	"github.com/gorilla/securecookie"
)

type repository interface {
	Register(buf bytes.Buffer) (*http.Cookie, error)
	UserLogin(buf bytes.Buffer) (*http.Cookie, error)
	GetUserFromCookie(cookieValue string) (string, error)
	GetServerSettings() *config.ServerSettings
	PostUserOrders(usr string, buf bytes.Buffer) error
	GetUserOrders(usr string) ([]byte, error)
	GetUserBalance(usr string) ([]byte, error)
	PostBalanceWithdraw(usr string, buf bytes.Buffer) error
	GetUserWithdrawals(usr string) ([]byte, error)
}

type Service struct {
	repo   repository
	secure *securecookie.SecureCookie
}

type LgnPsw struct {
	Lgn string `json:"login"`
	Psw string `json:"password"`
}

// Register implements server.svc.
func (s *Service) Register(buf bytes.Buffer) (*http.Cookie, error) {
	var err error
	lgnPsw := LgnPsw{}
	if err = json.Unmarshal(buf.Bytes(), &lgnPsw); err != nil {
		return nil, err
	}
	fmt.Println(lgnPsw)
	if err = dbstorage.InsertUser(s.repo.GetServerSettings().DB, lgnPsw.Lgn, lgnPsw.Psw); err != nil {
		return nil, err
	}
	//если новый пользователь, валидный JSON то получаем куку - зашифрованный юзер
	cookie, err := s.SetUserCookie(lgnPsw.Lgn)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

func (s *Service) UserLogin(buf bytes.Buffer) (*http.Cookie, error) {
	var err error
	lgnPsw := LgnPsw{}
	if err = json.Unmarshal(buf.Bytes(), &lgnPsw); err != nil {
		return nil, err
	}
	fmt.Println(lgnPsw)
	if err = dbstorage.CheckLgnPsw(s.repo.GetServerSettings().DB, lgnPsw.Lgn, lgnPsw.Psw); err != nil {
		return nil, err
	}
	//fmt.Println(err)
	//если новый пользователь, валидный JSON то получаем куку - зашифрованный юзер
	cookie, err := s.SetUserCookie(lgnPsw.Lgn)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

func (s *Service) GetAdresRun() string {
	return s.repo.GetServerSettings().AdresRun
}

/*func (s *Service) GetServerSettings() *config.ServerSettings {
	return &s.serverSettings
}*/

func NewService() (*Service, error) {

	serverSettings, err := config.InitServerSettings()
	if err != nil {
		log.Println(err)
	}
	//service.serverSettings = *serverSettings
	service := Service{repo: serverSettings}
	//fmt.Println(serverSettings.AdresBase)
	serverSettings.DB, err = sql.Open("pgx", serverSettings.AdresBase)
	if err != nil {
		return nil, err
	}

	if err := dbstorage.CreateTables(serverSettings.DB); err != nil {
		return nil, err
	}
	service.InitSecure()
	return &service, nil
}

func (s *Service) InitSecure() {
	var hashKey = []byte(s.repo.GetServerSettings().SecretCookies)
	var blockKey = []byte("a-lot-secret-qwe")
	s.secure = securecookie.New(hashKey, blockKey)
}
