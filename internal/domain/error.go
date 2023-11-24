package domain

import "errors"

var ErrNotEnoughFunds = errors.New("not enough funds")
var ErrNotFound = errors.New("record not found")
var ErrDuplicateOrderUploadSameUser = errors.New("ErrDuplicateOrderUploadSameUser")
var ErrDuplicateOrderUploadOtherUser = errors.New("ErrDuplicateOrderUploadOtherUser")
