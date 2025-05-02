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
	id        int32
	roleType  int32
	createdAt time.Time
}

type submission struct {
	id             int32
	retryCount     int32
	state          int32
	stateUpdatedAt *time.Time
}
