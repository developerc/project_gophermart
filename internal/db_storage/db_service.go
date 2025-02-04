package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/developerc/project_gophermart/internal/general"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type ErrorLgnPsw struct {
	s string
}

func (e *ErrorLgnPsw) Error() string {
	return e.s
}

func (e *ErrorLgnPsw) AsLgnPswWrong(err error) bool {
	return errors.As(err, &e)
}

//-----

func CreateTables(db *sql.DB) error {
	const duration uint = 20
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
	defer cancel()
	const crusr string = "CREATE TABLE IF NOT EXISTS usr_table( uuid serial primary key, " +
		"usr TEXT CONSTRAINT must_be_different_usr UNIQUE, psw TEXT)"
	_, err := db.ExecContext(ctx, crusr)
	if err != nil {
		return err
	}

	const crord string = "CREATE TABLE IF NOT EXISTS orders_table( uuid serial primary key, " +
		"usr TEXT, order_numb TEXT, status TEXT NOT NULL DEFAULT 'NEW', accrual REAL NOT NULL DEFAULT 0.0, date_time TIMESTAMP NOT NULL DEFAULT NOW())"
	_, err = db.ExecContext(ctx, crord)
	if err != nil {
		return err
	}

	const crwd string = "CREATE TABLE IF NOT EXISTS withdraw_table( uuid serial primary key, " +
		"usr TEXT, order_numb TEXT, withdraw REAL NOT NULL DEFAULT 0.0, date_time TIMESTAMP NOT NULL DEFAULT NOW())"
	_, err = db.ExecContext(ctx, crwd)
	if err != nil {
		return err
	}

	const pgcrpt string = "CREATE EXTENSION pgcrypto"
	_, err = db.ExecContext(ctx, pgcrpt)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func InsertUser(db *sql.DB, usr, psw string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx, "INSERT INTO usr_table (usr, psw) values ($1, crypt($2, gen_salt('md5')))", usr, psw)
	if err != nil {
		return err
	}
	return nil
}

func CheckLgnPsw(db *sql.DB, usr, psw string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, "SELECT (psw = crypt($2, psw)) AS password_match FROM usr_table WHERE usr = $1 ", usr, psw)
	if err != nil {
		return err
	}
	defer rows.Close()

	cntrRows := 0
	for rows.Next() {
		cntrRows++
		var password_match bool
		err = rows.Scan(&password_match)
		//fmt.Println("password_match: ", password_match)
		if err != nil {
			return err
		}
		if !password_match {
			//return errors.New("login or password is not valid")
			return &ErrorLgnPsw{"login or password is not valid"}
		}
	}
	if cntrRows == 0 {
		//return errors.New("login or password is not valid")
		return &ErrorLgnPsw{"login or password is not valid"}
	}

	return nil
}

func UploadOrder(db *sql.DB, usr, orderNum string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	rows, err := db.QueryContext(ctx, "SELECT usr FROM orders_table WHERE order_numb = $1 ", orderNum)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	cntrRows := 0
	var usrInTable string
	for rows.Next() {
		cntrRows++
		err = rows.Scan(&usrInTable)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if cntrRows > 0 {
		if usrInTable == usr {
			tx.Rollback()
			//return errors.New("alredy exists order the same usr 200") //добавить типизированную ошибку
			return &general.ErrorExistsOrderSame{}
		} else {
			tx.Rollback()
			//return errors.New("alredy exists order other usr 409")
			return &general.ErrorExistsOrderOther{}
		}
	}

	_, err = db.ExecContext(ctx, "INSERT INTO orders_table (usr, order_numb) values ($1, $2)", usr, orderNum)
	if err != nil {
		return err
	}

	// завершаем транзакцию
	return tx.Commit()
}

func GetUserOrders(db *sql.DB, usr string) ([]general.UploadedOrder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, "SELECT order_numb, status, accrual, date_time from orders_table WHERE usr = $1 ORDER BY date_time DESC", usr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var arrUploadedOrder []general.UploadedOrder = make([]general.UploadedOrder, 0)
	for rows.Next() {
		uploadedOrder := general.UploadedOrder{}
		var number string
		var status string
		var accrual float64
		var uploaded_at time.Time
		err = rows.Scan(&number, &status, &accrual, &uploaded_at)
		if err != nil {
			return nil, err
		}
		uploadedOrder.Number = number
		uploadedOrder.Status = status
		uploadedOrder.Accrual = accrual
		//fmt.Println(uploaded_at.Format("2006-01-02T15:04:05-07:00"))
		/*uploaded_at, err = time.Parse(time.RFC3339, uploaded_at.String())
		if err != nil {
			return nil, err
		}*/
		//fmt.Println(uploaded_at)
		uploadedOrder.UploadedAt = uploaded_at.Format("2006-01-02T15:04:05-07:00")

		arrUploadedOrder = append(arrUploadedOrder, uploadedOrder)
	}

	return arrUploadedOrder, nil
}

func GetUserBalance(db *sql.DB, usr string) (general.UserBalance, error) {
	var sumAccrual float64
	var sumWithdraw float64
	userBalance := general.UserBalance{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	tx, err := db.Begin()
	if err != nil {
		return userBalance, err
	}
	rows, err := db.QueryContext(ctx, "SELECT COALESCE(SUM(accrual), 0 ) from orders_table WHERE usr = $1", usr)
	if err != nil {
		tx.Rollback()
		return userBalance, err
	}
	defer rows.Close()
	cntrRows := 0
	for rows.Next() {
		cntrRows++
		err = rows.Scan(&sumAccrual)
		if err != nil {
			tx.Rollback()
			return userBalance, err
		}
	}
	if cntrRows == 0 {
		sumAccrual = 0
	}
	userBalance.Current = sumAccrual

	rows2, err := db.QueryContext(ctx, "SELECT COALESCE(SUM(withdraw), 0 ) from withdraw_table WHERE usr = $1", usr)
	if err != nil {
		tx.Rollback()
		return userBalance, err
	}
	defer rows2.Close()
	cntrRows = 0
	for rows2.Next() {
		cntrRows++
		err = rows2.Scan(&sumWithdraw)
		if err != nil {
			tx.Rollback()
			return userBalance, err
		}
	}
	if cntrRows == 0 {
		sumWithdraw = 0
	}
	userBalance.Withdrawn = sumWithdraw
	err = tx.Commit()
	if err != nil {
		return userBalance, err
	}
	return userBalance, nil
}

func CheckUsrOrderNumb(db *sql.DB, usr string, order string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, "SELECT COUNT(*) from orders_table WHERE usr = $1 AND order_numb = $2", usr, order)
	if err != nil {
		return err
	}
	var cnt int
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&cnt)
		if err != nil {
			return err
		}
		fmt.Println("cnt:", cnt)
	}
	if cnt != 1 {
		//return errors.New("invalid number of order")
		return &general.ErrorNumOrder{}
	}
	return nil
}

func BalanceWithdraw(db *sql.DB, usr string, order string, sum float64) error {
	var sumAccrual float64
	var sumWithdraw float64
	var diffSum float64
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, "SELECT COALESCE(SUM(accrual), 0 ) from orders_table WHERE usr = $1", usr)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&sumAccrual)
		if err != nil {
			return err
		}
	}
	fmt.Println("sumAccrual:", sumAccrual)
	rows2, err := db.QueryContext(ctx, "SELECT COALESCE(SUM(withdraw), 0 ) from withdraw_table WHERE usr = $1", usr)
	if err != nil {
		return err
	}
	defer rows2.Close()
	for rows2.Next() {
		err = rows2.Scan(&sumWithdraw)
		if err != nil {
			return err
		}
	}
	fmt.Println("sumWithdraw:", sumWithdraw)
	diffSum = sumAccrual - sumWithdraw
	fmt.Println("diffSum:", diffSum)
	if diffSum < sum {
		//return errors.New("not enough of loyalty points")
		return &general.ErrorLoyaltyPoints{}
	}

	_, err = db.ExecContext(ctx, "INSERT INTO withdraw_table (usr, order_numb, withdraw) values ($1, $2, $3)", usr, order, sum)
	if err != nil {
		return err
	}
	return nil
}
