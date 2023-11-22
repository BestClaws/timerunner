package main

import (
	"fmt"
	"strings"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	dur := (string(b))[1 : len(b)-1]
	duration, err := time.ParseDuration(dur)
	if err != nil {
		return fmt.Errorf("err invalid duration: %s", err)
	}

	d.Duration = duration
	return nil
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	var str strings.Builder
	str.WriteByte('"')
	h := int(d.Duration.Hours()) % 24
	m := int(d.Duration.Minutes()) % 60
	s := int(d.Duration.Seconds()) % 60
	if h > 0 {
		str.WriteString(fmt.Sprintf("%dh", h))
	}
	if m > 0 {
		str.WriteString(fmt.Sprintf("%dm", m))
	}
	if s > 0 {
		str.WriteString(fmt.Sprintf("%ds", s))
	}
	str.WriteByte('"')
	return []byte(str.String()), nil
}

// TwelveAmTime resets hours, mins, seconds to zero for the given time
func TwelveAmTime(t time.Time) time.Time {
	year, month, date := t.Date()
	return time.Date(year, month, date, 0, 0, 0, 0, t.Location())
}
