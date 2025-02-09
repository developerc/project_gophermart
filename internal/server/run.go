package server

import (
	"log"
	"net/http"

	"github.com/developerc/project_gophermart/internal/service"
)

func Run() error {
	service, err := service.NewService()
	if err != nil {
		return err
	}
	server, err := NewServer(service)
	if err != nil {
		return err
	}
	routes := server.SetupRoutes()
	err = http.ListenAndServe(service.GetAdresRun(), routes)
	if err != nil {
		log.Println(err)
	}

	return nil
}
