package loyalty

import (
	"database/sql"
	"log"
	"time"

	dbstorage "github.com/developerc/project_gophermart/internal/db_storage"
)

func RunLoyalty(db *sql.DB, adresAccrual string) {
	go func() {
		for {
			chanCnt := 5
			//из таблицы order_table выбираем все номера заказов где status равен NEW или REGISTERED или PROCESSING
			arrOrderNumb, err := dbstorage.GetOrderNumbs(db)
			if err != nil {
				log.Println(err)
				continue
			}
			//fmt.Println("chanCnt:", chanCnt, ", arrOrderNumb:", arrOrderNumb)
			DoRequests(db, chanCnt, arrOrderNumb, adresAccrual)

			//будем проверять переменную, в которую может записаться необходимая задержка
			time.Sleep(1 * time.Second)
		}
	}()
}
