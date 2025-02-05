package service

import (
	"bytes"
	"encoding/json"

	//"fmt"

	//"errors"

	"net/http"

	dbstorage "github.com/developerc/project_gophermart/internal/db_storage"
	"github.com/developerc/project_gophermart/internal/general"
)

type User struct {
	Name string
}

type OrderSum struct {
	Order string
	Sum   float64
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
	//fmt.Println("from PostUserOrders usr", usr)
	//var numOrderStr string
	//numOrderStr := buf.String()
	//numOrderBytes := buf.Bytes()
	//fmt.Println(numOrderStr)
	//fmt.Println(numOrderBytes)
	//fmt.Println("len(numOrderStr):", len(numOrderStr), ", len(numOrderBytes):", len(numOrderBytes))
	//делать не будем, хендлера на удаление юзера нет! проверка существует ли до сих пор такой юзер в таблице (мог быть раньше зарегистрирован, потом удален)

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
	if err := dbstorage.UploadOrder(s.repo.GetServerSettings().DB, usr, buf.String()); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUserOrders(usr string) ([]byte, error) {
	//fmt.Println("from GetUserOrders usr", usr)
	//var arrUploadedOrder []general.UploadedOrder
	arrUploadedOrder, err := dbstorage.GetUserOrders(s.repo.GetServerSettings().DB, usr)
	if err != nil {
		return nil, err
	}
	//fmt.Println(arrUploadedOrder)
	jsonBytes, err := json.Marshal(arrUploadedOrder)
	if err != nil {
		return nil, err
	}
	if len(arrUploadedOrder) == 0 {
		//return nil, errors.New("no data 204")
		return nil, &general.ErrorNoContent{}
	}

	return jsonBytes, nil
}

func (s Service) GetUserBalance(usr string) ([]byte, error) {
	//fmt.Println("from GetUserBalance usr", usr)
	userBalance, err := dbstorage.GetUserBalance(s.repo.GetServerSettings().DB, usr)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(userBalance)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func (s *Service) PostBalanceWithdraw(usr string, buf bytes.Buffer) error {
	var err error
	orderSum := OrderSum{}
	if err = json.Unmarshal(buf.Bytes(), &orderSum); err != nil {
		return err
	}
	//fmt.Println(orderSum)
	err = dbstorage.CheckUsrOrderNumb(s.repo.GetServerSettings().DB, usr, orderSum.Order)
	//fmt.Println("err1:", err)
	if err != nil {
		return err
	}
	err = dbstorage.BalanceWithdraw(s.repo.GetServerSettings().DB, usr, orderSum.Order, orderSum.Sum)
	//fmt.Println("err2:", err)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUserWithdrawals(usr string) ([]byte, error) {
	arrWithdrawOrder, err := dbstorage.GetUserWithdrawals(s.repo.GetServerSettings().DB, usr)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(arrWithdrawOrder)
	if err != nil {
		return nil, err
	}
	if len(arrWithdrawOrder) == 0 {
		//return nil, errors.New("no data 204")
		return nil, &general.ErrorNoContent{}
	}

	return jsonBytes, nil
}
