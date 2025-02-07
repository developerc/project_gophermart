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
type ErrorLoyaltyPoints struct {
	s string
}

func (e *ErrorLoyaltyPoints) Error() string {
	return e.s
}

func (e *ErrorLoyaltyPoints) AsErrorNoContent(err error) bool {
	return errors.As(err, &e)
}

// ---
type UploadedOrder struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

type UserBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawOrder struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type LoyaltyOrder struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
