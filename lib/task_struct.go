package lib

import "time"

type Task struct {
	Email    string
	Deadline time.Time
	Title     string
}