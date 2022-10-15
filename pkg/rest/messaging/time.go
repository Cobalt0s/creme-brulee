package messaging

import (
	"fmt"
	"strings"
	"time"
)

const (
	timeNanoFormat = time.RFC3339Nano
	timeISOFormat  = "2006-01-02T15:04:05.000Z"
)

var (
	nilTime = (time.Time{}).UnixNano()
)

func NewNano3339Time(t time.Time) TimeNano3339 {
	return TimeNano3339{Time: t}
}

func ParseNano3339Time(s string) (*time.Time, error) {
	parse, err := time.Parse(timeNanoFormat, s)
	if err != nil {
		return nil, err
	}
	return &parse, nil
}

type TimeNano3339 struct {
	time.Time
}

func (ct *TimeNano3339) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(timeNanoFormat, s)
	return
}

func (ct *TimeNano3339) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(timeNanoFormat))), nil
}

func NewISO8601Time(t time.Time) TimeISO8601 {
	return TimeISO8601{Time: t}
}

func ParseISO8601Time(s string) (*time.Time, error) {
	parse, err := time.Parse(timeISOFormat, s)
	if err != nil {
		return nil, err
	}
	return &parse, nil
}

type TimeISO8601 struct {
	time.Time
}

func (ct *TimeISO8601) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(timeISOFormat, s)
	return
}

func (ct *TimeISO8601) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(timeISOFormat))), nil
}
