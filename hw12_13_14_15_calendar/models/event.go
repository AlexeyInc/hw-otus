package models

import "time"

type Event struct {
	ID           int64
	Title        string
	StartEvent   time.Time
	EndEvent     time.Time
	Description  string
	IDUser       int64
	Notification time.Time
}
