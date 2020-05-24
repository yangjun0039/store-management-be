package network

import (
	"time"
)

type Requester struct {
	UUID      string
	User      UserInfo
	timestamp time.Time
}

type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}
