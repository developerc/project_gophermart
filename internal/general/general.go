package general

import (
	"errors"
	//"time"
)

type ErrorNumOrder struct {
	s string
}

func (e *ErrorNumOrder) Error() string {
	return e.s
}

func (e *ErrorNumOrder) AsNumOrderWrong(err error) bool {
	return errors.As(err, &e)
}

// ---
type ErrorExistsOrderSame struct {
	s string
}

func (e *ErrorExistsOrderSame) Error() string {
	return e.s
}

func (e *ErrorExistsOrderSame) AsExistsOrderSame(err error) bool {
	return errors.As(err, &e)
}

// ---
type ErrorExistsOrderOther struct {
	s string
}

func (e *ErrorExistsOrderOther) Error() string {
	return e.s
}

func (e *ErrorExistsOrderOther) AsExistsOrderOther(err error) bool {
	return errors.As(err, &e)
}

// ---
type ErrorNoContent struct {
	s string
}

func (e *ErrorNoContent) Error() string {
	return e.s
}

func (e *ErrorNoContent) AsErrorNoContent(err error) bool {
	return errors.As(err, &e)
}

// ---
type UploadedOrder struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    int    `json:"accrual,omitempty"`
	UploadedAt string `json:"uploaded_at"`
}

type UserBalance struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}
