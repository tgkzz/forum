package model

import "time"

type Session struct {
	Id      int
	Token   string
	ExpTime time.Time
	UserId  int
}
