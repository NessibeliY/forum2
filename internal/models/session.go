package models

import "time"

type Session struct {
	SessionId  string
	UserId     string
	Token      string
	ExpireTime time.Time
}
