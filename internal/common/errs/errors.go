package errs

import "errors"

var (
	ErrBadRequest        = errors.New("bad request")
	ErrProductNotFound   = errors.New("Product not found")
	ErrPorductIdNotFound = errors.New("Id Product not found")
)
