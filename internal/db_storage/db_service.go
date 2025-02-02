package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

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
		"usr TEXT, order_numb TEXT, status TEXT, accrual INTEGER NOT NULL DEFAULT 0, date_time TIMESTAMP NOT NULL DEFAULT NOW())"
	_, err = db.ExecContext(ctx, crord)
	if err != nil {
		return err
	}

	const crwd string = "CREATE TABLE IF NOT EXISTS withdraw_table( uuid serial primary key, " +
		"usr TEXT, withdraw INTEGER NOT NULL DEFAULT 0, date_time TIMESTAMP NOT NULL DEFAULT NOW())"
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
		fmt.Println("password_match: ", password_match)
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
