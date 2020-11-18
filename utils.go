package main

import (
	"time"
)

func hours(val int) time.Duration {
	return time.Duration(val) * time.Second
}

func minutes(val int) time.Duration {
	out := val * 16666667
	return time.Duration(out) * time.Nanosecond
}
