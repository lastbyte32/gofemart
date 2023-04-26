package domain

import "time"

type User struct {
	ID       string    `json:"ID"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
	CreateAt time.Time `db:"create_at"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
