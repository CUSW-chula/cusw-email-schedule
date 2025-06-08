package lib

import (
	"database/sql"
	"log"
)

func QueryTasks(db *sql.DB) []Task {
	query := `
		SELECT
  			users.email,
  			tasks."endDate" AS deadline,
  			tasks.title
		FROM
  			task_assignments
		JOIN
  			users ON task_assignments."userId" = users.id
		JOIN
  			tasks ON task_assignments."taskId" = tasks.id
		WHERE
  			tasks."endDate"::date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '1 day'
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
		if err := rows.Scan(&t.Email, &t.Deadline, &t.Title); err != nil {
			log.Printf("Row scan failed: %v", err)
			continue
		}
		tasks = append(tasks, t)
	}

	return tasks
}