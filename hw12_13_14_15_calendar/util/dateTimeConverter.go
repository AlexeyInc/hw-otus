package util

import "time"

type Period string

func (s Period) GetTimePeriod() time.Time {
	switch s {
	case "Day":
		return time.Now().UTC()
	case "Week":
		return time.Now().AddDate(0, 0, 7).UTC()
	case "Month":
		return time.Now().AddDate(0, 0, 30).UTC()
	}
	return time.Now().UTC()
}
