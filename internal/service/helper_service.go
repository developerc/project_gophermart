package service

import (
	"bytes"
	"encoding/json"
	"log"

	//"fmt"

	//"errors"

	"net/http"

	dbstorage "github.com/developerc/project_gophermart/internal/db_storage"
	"github.com/developerc/project_gophermart/internal/general"
	"github.com/theplant/luhn"
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
	//делать не будем, хендлера на удаление юзера нет! проверка существует ли до сих пор такой юзер в таблице (мог быть раньше зарегистрирован, потом удален)
	log.Println("from hs PostUserOrders usr:", usr, ", order:", buf.String())
	//проверим номер заказа на Луна. Номер заказа может не быть в таблице заказов
	err := checkLuhna(buf.String())
	log.Println("from hs PostUserOrders, checkLuhna: ", err)
	if err != nil {
		return err
	}
	//проверка валидности строки запроса
	/*var intNum int = 0
	for _, runeValue := range buf.String() {
		//var isdig bool = false
		if runeValue < 48 || runeValue > 57 {
			//return errors.New("invalid numer of order")
			return &general.ErrorNumOrder{}
		}
		//проверка на Луну
		intNum = intNum*10 + int(runeValue-48)
	}
	if !luhn.Valid(intNum) {
		return &general.ErrorNumOrder{}
	}*/
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
		log.Println("from hs GetUserOrders usr:", usr, ", Error dbstorage:", err)
		return nil, err
	}
	log.Println("from hs GetUserOrders usr:", usr, ", arrUploadedOrder:", arrUploadedOrder)
	jsonBytes, err := json.Marshal(arrUploadedOrder)
	if err != nil {
		log.Println("from hs GetUserOrders usr:", usr, ", Error Marshal:", err)
		return nil, err
	}
	if len(arrUploadedOrder) == 0 {
		//return nil, errors.New("no data 204")
		log.Println("from hs GetUserOrders usr:", usr, ", ErrorNoContent")
		return nil, &general.ErrorNoContent{}
	}

	return jsonBytes, nil
}

func (s Service) GetUserBalance(usr string) ([]byte, error) {
	//fmt.Println("from GetUserBalance usr", usr)
	userBalance, err := dbstorage.GetUserBalance2(s.repo.GetServerSettings().DB, usr)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(userBalance)
	if err != nil {
		return nil, err
	}
	log.Println("from GetUserBalance usr: ", usr, ", userBalance: ", userBalance)
	return jsonBytes, nil
}

func (s *Service) PostBalanceWithdraw(usr string, buf bytes.Buffer) error {
	var err error
	orderSum := OrderSum{}
	if err = json.Unmarshal(buf.Bytes(), &orderSum); err != nil {
		return err
	}
	log.Println("from hs PostBalanceWithdraw, usr: ", usr, ", orderSum: ", orderSum)
	//err = dbstorage.CheckUsrOrderNumb(s.repo.GetServerSettings().DB, usr, orderSum.Order)
	//проверим номер заказа на Луна. Номер заказа может не быть в таблице заказов
	err = checkLuhna(orderSum.Order)
	log.Println("from hs PostBalanceWithdraw, checkLuhna: ", err)
	if err != nil {
		return err
	}
	err = dbstorage.BalanceWithdraw2(s.repo.GetServerSettings().DB, usr, orderSum.Order, orderSum.Sum)
	log.Println("from hs PostBalanceWithdraw, BalanceWithdraw err: ", err)
	if err != nil {
		return err
	}
	return nil
}

func checkLuhna(order string) error {
	//проверка валидности строки запроса
	intNum := 0
	for _, runeValue := range order {
		if runeValue < 48 || runeValue > 57 {
			return &general.ErrorNumOrder{}
		}
		//проверка на Луну
		intNum = intNum*10 + int(runeValue-48)
	}
	if !luhn.Valid(intNum) {
		return &general.ErrorNumOrder{}
	}
	return nil
}

func (s *Service) GetUserWithdrawals(usr string) ([]byte, error) {
	arrWithdrawOrder, err := dbstorage.GetUserWithdrawals2(s.repo.GetServerSettings().DB, usr)
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
