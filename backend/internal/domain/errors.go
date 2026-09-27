package domain

import (
	"errors"
)

var ErrOrderNotFound = errors.New("order not found")
var ErrInvalidOrder = errors.New("invalid order")
var ErrOrderAlreadyExist = errors.New("order already exists")
