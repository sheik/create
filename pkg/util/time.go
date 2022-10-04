package util

import "time"

func CurrentTime() string {
	return time.Now().Format("2006-01-02T15:04:05.52Z")
}
