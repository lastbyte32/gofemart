package domain

import "time"

type User struct {
	ID       string    `json:"ID"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
	CreateAt time.Time `db:"create_at"`
}
