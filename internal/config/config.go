package config

import (
	"bytes"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
)

/*type Settings interface {
	GetAdresRun() string
}*/

type ServerSettings struct {
	//Settings
	AdresRun      string
	AdresBase     string
	AdresAccrual  string
	SecretCookies string
	DB            *sql.DB
}

func InitServerSettings() (*ServerSettings, error) {
	//var err error
	serverSettings := &ServerSettings{}
	ar := flag.String("a", "localhost:8081", "address running server")
	dbStorage := flag.String("d", "", "address connect to DB")
	asa := flag.String("r", "", "address accrual system")
	sec := flag.String("s", "super_secret", "secret for cookies")
	flag.Parse()

	val, ok := os.LookupEnv("RUN_ADDRESS")
	if !ok || val == "" {
		serverSettings.AdresRun = *ar
		log.Println("AdresRun from flag:", serverSettings.AdresRun)
		//serverSettings.Logger.Info("AdresRun from flag:", zap.String("address", serverSettings.AdresRun))
	} else {
		serverSettings.AdresRun = val
		log.Println("AdresRun from env:", serverSettings.AdresRun)
		//serverSettings.Logger.Info("AdresRun from env:", zap.String("address", serverSettings.AdresRun))
	}

	val, ok = os.LookupEnv("DATABASE_URI")
	if !ok || val == "" {
		serverSettings.AdresBase = *dbStorage
		log.Println("DbStorage from flag:", serverSettings.AdresBase)
		//if isFlagPassed("d") && (serverSettings.DBStorage != "") {
		//serverSettings.TypeStorage = DBStorage
		//}
		//serverSettings.Logger.Info("DbStorage from flag:", zap.String("storage", serverSettings.DBStorage))
	} else {
		//serverSettings.TypeStorage = DBStorage
		serverSettings.AdresBase = val
		log.Println("DbStorage from env:", serverSettings.AdresBase)
		//serverSettings.Logger.Info("DbStorage from env:", zap.String("storage", serverSettings.DBStorage))
	}

	val, ok = os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if !ok || val == "" {
		serverSettings.AdresAccrual = *asa
		log.Println("AdresAccrual from flag:", serverSettings.AdresAccrual)
		//serverSettings.Logger.Info("AdresRun from flag:", zap.String("address", serverSettings.AdresRun))
	} else {
		serverSettings.AdresAccrual = val
		log.Println("AdresAccrual from env:", serverSettings.AdresAccrual)
		//serverSettings.Logger.Info("AdresRun from env:", zap.String("address", serverSettings.AdresRun))
	}

	val, ok = os.LookupEnv("SECRET_FOR_COOKIES")
	if !ok || val == "" {
		serverSettings.SecretCookies = *sec
		log.Println("SecretCookies from flag:", serverSettings.SecretCookies)
		//serverSettings.Logger.Info("AdresRun from flag:", zap.String("address", serverSettings.AdresRun))
	} else {
		serverSettings.SecretCookies = val
		log.Println("SecretCookies from env:", serverSettings.SecretCookies)
		//serverSettings.Logger.Info("AdresRun from env:", zap.String("address", serverSettings.AdresRun))
	}
	return serverSettings, nil
}

func (s *ServerSettings) GetAdresRun() string {
	return s.AdresRun
}

func (s *ServerSettings) GetServerSettings() *ServerSettings {
	return s
}

func (s *ServerSettings) Register(buf bytes.Buffer) (*http.Cookie, error) {
	return nil, nil
}

func (s *ServerSettings) UserLogin(buf bytes.Buffer) (*http.Cookie, error) {
	return nil, nil
}

func (s *ServerSettings) GetUserFromCookie(cookieValue string) (string, error) {
	return "", nil
}

func (s *ServerSettings) PostUserOrders(usr string, buf bytes.Buffer) error {
	return nil
}

func (s *ServerSettings) GetUserOrders(usr string) ([]byte, error) {
	return nil, nil
}

func (s *ServerSettings) GetUserBalance(usr string) ([]byte, error) {
	return nil, nil
}
