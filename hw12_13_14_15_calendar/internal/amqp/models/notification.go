package contracts

import "time"

type Notification struct {
	IdEvent    int64
	EventTitle string
	EventStart time.Time
	IdUser     int64
}
