package database

import (
	"time"
)

type user struct {
	id        int32
	username  string
	password  string
	roleId    int32
	createdAt time.Time
}

type role struct {
	id        string
	roleType  int32
	createdAt time.Time
}
