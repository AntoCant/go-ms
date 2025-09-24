package errors

import "errors"

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrInvalidProductID = errors.New("invalid product id")
)
