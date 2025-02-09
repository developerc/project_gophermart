package loyalty

import (
	"database/sql"
	"log"
	"time"

	dbstorage "github.com/developerc/project_gophermart/internal/db_storage"
)

var retryAfterSec int

func RunLoyalty(db *sql.DB, adresAccrual string) {
	go func() {
		for {
			chanCnt := 5
			arrOrderNumb, err := dbstorage.GetOrderNumbs(db)
			if err != nil {
				log.Println(err)
				continue
			}
			DoRequests(db, chanCnt, arrOrderNumb, adresAccrual)
			time.Sleep(time.Duration(retryAfterSec) * time.Second)
			retryAfterSec = 0
			time.Sleep(1 * time.Second)
		}
	}()
}
