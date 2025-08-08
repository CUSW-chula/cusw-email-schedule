package lib

import (
	"database/sql"
	"fmt"
	"log"
)

func QueryTasks(db *sql.DB) []Task {
	query := `
		SELECT
			tasks.id,
			tasks.title,
			tasks.description,
			tasks.status,
			tasks."projectId",
			tasks."startDate",
			tasks."endDate",
			tasks.budget,
			projects.title AS project_title,
			assignee.name AS assignee_name,
			assignee.email AS assignee_email,
			COALESCE(creator.name, 'System') AS assignor_name
		FROM
			task_assignments
		JOIN
			users assignee ON task_assignments."userId" = assignee.id
		JOIN
			tasks ON task_assignments."taskId" = tasks.id
		JOIN
			projects ON tasks."projectId" = projects.id
		LEFT JOIN
			users creator ON tasks."createdById" = creator.id
		WHERE
			tasks."endDate"::date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '1 day'
			AND assignee.activated = true
			AND tasks.status != 'Done'
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
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.ProjectID,
			&t.StartDate,
			&t.EndDate,
			&t.Budget,
			&t.ProjectTitle,
			&t.AssigneeName,
			&t.AssigneeEmail,
			&t.AssignorName,
		); err != nil {
			log.Printf("Row scan failed: %v", err)
			continue
		}
		tasks = append(tasks, t)
	}

	return tasks
}

// QueryTasksDueInDays ดึงงานที่กำลังจะครบกำหนดในจำนวนวันที่กำหนด
func QueryTasksDueInDays(db *sql.DB, days int) []Task {
	query := `
		SELECT
			tasks.id,
			tasks.title,
			tasks.description,
			tasks.status,
			tasks."projectId",
			tasks."startDate",
			tasks."endDate",
			tasks.budget,
			projects.title AS project_title,
			assignee.name AS assignee_name,
			assignee.email AS assignee_email,
			creator.name AS assignor_name
		FROM
			task_assignments
		JOIN
			users assignee ON task_assignments."userId" = assignee.id
		JOIN
			tasks ON task_assignments."taskId" = tasks.id
		JOIN
			projects ON tasks."projectId" = projects.id
		LEFT JOIN
			users creator ON tasks."createdById" = creator.id
		WHERE
			tasks."endDate"::date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '%d days'
			AND assignee.activated = true
			AND tasks.status != 'Done'
	`

	formattedQuery := fmt.Sprintf(query, days)
	rows, err := db.Query(formattedQuery)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return nil
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.ProjectID,
			&t.StartDate,
			&t.EndDate,
			&t.Budget,
			&t.ProjectTitle,
			&t.AssigneeName,
			&t.AssigneeEmail,
			&t.AssignorName,
		); err != nil {
			log.Printf("Row scan failed: %v", err)
			continue
		}
		tasks = append(tasks, t)
	}

	return tasks
}

// QueryOverdueTasks ดึงงานที่เลยกำหนดแล้ว
func QueryOverdueTasks(db *sql.DB) []Task {
	query := `
		SELECT
			tasks.id,
			tasks.title,
			tasks.description,
			tasks.status,
			tasks."projectId",
			tasks."startDate",
			tasks."endDate",
			tasks.budget,
			projects.title AS project_title,
			assignee.name AS assignee_name,
			assignee.email AS assignee_email,
			creator.name AS assignor_name
		FROM
			task_assignments
		JOIN
			users assignee ON task_assignments."userId" = assignee.id
		JOIN
			tasks ON task_assignments."taskId" = tasks.id
		JOIN
			projects ON tasks."projectId" = projects.id
		LEFT JOIN
			users creator ON tasks."createdById" = creator.id
		WHERE
			tasks."endDate" < CURRENT_TIMESTAMP
			AND assignee.activated = true
			AND tasks.status != 'Done'
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
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.ProjectID,
			&t.StartDate,
			&t.EndDate,
			&t.Budget,
			&t.ProjectTitle,
			&t.AssigneeName,
			&t.AssigneeEmail,
			&t.AssignorName,
		); err != nil {
			log.Printf("Row scan failed: %v", err)
			continue
		}
		tasks = append(tasks, t)
	}

	return tasks
}
