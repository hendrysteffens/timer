package util

import "time"

const (
	typed  = iota
	server = iota
)

type dayTimes struct {
	typeTime int
	time     time.Time
}

func Add(dayTime time.Time) {

}
