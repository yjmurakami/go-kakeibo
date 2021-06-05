package clock

import "time"

type Clock interface {
	Now() time.Time
	Today() time.Time
}

type SystemClock struct{}

func (s SystemClock) Now() time.Time {
	return time.Now()
}

func (s SystemClock) Today() time.Time {
	now := s.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	return today
}
