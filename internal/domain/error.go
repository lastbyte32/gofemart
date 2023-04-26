package domain

import "errors"

var ErrNotEnoughFunds = errors.New("not enough funds")
var ErrNotFound = errors.New("record not found")
