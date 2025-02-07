package loyalty

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	dbstorage "github.com/developerc/project_gophermart/internal/db_storage"
	"github.com/developerc/project_gophermart/internal/general"
)

// Semaphore структура семафора
type Semaphore struct {
	semaCh chan struct{}
}

// NewSemaphore создает семафор с буферизованным каналом емкостью maxReq
func NewSemaphore(maxReq int) *Semaphore {
	return &Semaphore{
		semaCh: make(chan struct{}, maxReq),
	}
}

// когда горутина запускается, отправляем пустую структуру в канал semaCh
func (s *Semaphore) Acquire() {
	s.semaCh <- struct{}{}
}

// когда горутина завершается, из канала semaCh убирается пустая структура
func (s *Semaphore) Release() {
	<-s.semaCh
}

func DoRequests(db *sql.DB, chanCnt int, arrOrderNumb []int, adresAccrual string) {
	var wg sync.WaitGroup
	semaphore := NewSemaphore(chanCnt)
	// создаем len(arrOrderNumb) горутин
	for idx := 0; idx < len(arrOrderNumb); idx++ {
		wg.Add(1)
		go func(orderNumb int) {
			semaphore.Acquire()
			defer wg.Done()
			defer semaphore.Release()
			//отправляем GET запрос в СРБНЛ, ответ записываем в order_table
			if err := ReqLoyalty(db, adresAccrual, orderNumb); err != nil {
				log.Println(err)
			}
		}(arrOrderNumb[idx])
	}
	wg.Wait()
}

func ReqLoyalty(db *sql.DB, adresAccrual string, orderNumb int) error {
	response, err := http.Get(adresAccrual + "/api/orders/" + strconv.FormatInt(int64(orderNumb), 10))
	if err != nil {
		log.Println(err)
		return err
	}
	//log.Println(response.StatusCode)
	stCode := response.StatusCode
	if stCode == 204 {
		return errors.New("response Status code: 204")
	}
	if stCode == 500 {
		return errors.New("response Status code: 500")
	}
	if stCode == 429 {
		//обработаем Body, приостановим запросы
		return errors.New("response Status code: 429")
	}
	//если статускод 200
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	//
	//orderSum := OrderSum{}
	loyaltyOrder := general.LoyaltyOrder{}
	if err = json.Unmarshal(body, &loyaltyOrder); err != nil {
		return err
	}
	//log.Println(loyaltyOrder)
	//делать Update заказа
	err = dbstorage.SetStatusAccrual(db, loyaltyOrder.Order, loyaltyOrder.Status, loyaltyOrder.Accrual)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
