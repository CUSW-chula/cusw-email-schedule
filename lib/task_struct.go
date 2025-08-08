package lib

import "time"

type Task struct {
	ID            string
	Title         string
	Description   string
	Status        string
	ProjectID     string
	StartDate     *time.Time
	EndDate       *time.Time
	Budget        float64
	ProjectTitle  string
	AssigneeName  string
	AssignorName  string
	AssigneeEmail string
}
