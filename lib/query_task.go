package lib

import (
	"database/sql"
	"log"
)

func QueryTasks(db *sql.DB) []Task {
	query := `
		SELECT email, deadline, task 
		FROM tasks 
		WHERE deadline >= CURRENT_DATE
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return nil
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.Email, &t.Deadline, &t.Task); err != nil {
			log.Printf("Row scan failed: %v", err)
			continue
		}
		tasks = append(tasks, t)
	}

	return tasks
}