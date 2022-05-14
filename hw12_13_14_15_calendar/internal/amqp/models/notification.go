package contracts

import "time"

type Notification struct {
	IDEvent    int64
	EventTitle string
	EventStart time.Time
	IDUser     int64
}
