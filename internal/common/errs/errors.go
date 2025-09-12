package errs

import "errors"

var (
	ErrBadRequest = errors.New("bad request")
)

var (
	PorductIdNotFound = errors.New("Id Product not found")
)

var (
	PorductNotFound = errors.New("Product not found")
)
