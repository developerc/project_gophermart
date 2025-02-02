package general

import "errors"

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
