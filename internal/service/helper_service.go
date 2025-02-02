package service

import (
	"bytes"
	"fmt"
	"net/http"

	dbstorage "github.com/developerc/project_gophermart/internal/db_storage"
	"github.com/developerc/project_gophermart/internal/general"
)

type User struct {
	Name string
}

func (s *Service) SetUserCookie(usr string) (*http.Cookie, error) {
	//var usr string
	var cookie *http.Cookie
	//var err error
	u := &User{
		Name: usr,
	}

	//if len(cookieValue) == 0 {
	u.Name = usr
	if encoded, err := s.secure.Encode("user", u); err == nil {
		cookie = &http.Cookie{
			Name:  "user",
			Value: encoded,
		}
		return cookie, nil
	} else {
		return nil, err
	}
	//}
	//return nil,  nil
}

func (s *Service) GetUserFromCookie(cookieValue string) (string, error) {
	var usr string
	u := &User{
		Name: usr,
	}
	if err := s.secure.Decode("user", cookieValue, u); err != nil {
		return "", err
	}
	//fmt.Println("u: ", u)

	return u.Name, nil
}

func (s *Service) PostUserOrders(usr string, buf bytes.Buffer) error {
	fmt.Println("from PostUserOrders usr", usr)
	//var numOrderStr string
	//numOrderStr := buf.String()
	//numOrderBytes := buf.Bytes()
	//fmt.Println(numOrderStr)
	//fmt.Println(numOrderBytes)
	//fmt.Println("len(numOrderStr):", len(numOrderStr), ", len(numOrderBytes):", len(numOrderBytes))
	//проверка валидности строки запроса
	for _, runeValue := range buf.String() {
		//var isdig bool = false
		if runeValue < 48 || runeValue > 57 {
			//return errors.New("invalid numer of order")
			return &general.ErrorNumOrder{}
		}
		//fmt.Println(runeValue, isdig)
	}
	//Загружаем заказ. проверка номер заказа уже загружен? кем?
	if err := dbstorage.LoadOrder(s.repo.GetServerSettings().DB, usr, buf.String()); err != nil {
		return err
	}
	return nil
}
