package lib

import (
	"database/sql"
	"log"
)

// QueryTasksByUserEmail ดึงงานของผู้ใช้เฉพาะคนจาก email สำหรับ 1 วันข้างหน้า
func QueryTasksByUserEmail(db *sql.DB, email string) []Task {
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
			assignee.email = $1
			AND assignee.activated = true
			AND tasks."endDate"::date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '1 day'
		ORDER BY tasks."endDate" ASC
	`

	rows, err := db.Query(query, email)
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

// GetUniqueUserEmails ดึง email ของผู้ใช้ที่มีงานใกล้ครบกำหนด
func GetUniqueUserEmails(db *sql.DB) []string {
	query := `
		SELECT DISTINCT assignee.email
		FROM task_assignments
		JOIN users assignee ON task_assignments."userId" = assignee.id
		JOIN tasks ON task_assignments."taskId" = tasks.id
		WHERE tasks."endDate"::date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '1 day'
		AND assignee.activated = true
		AND tasks.status != 'Done'
		ORDER BY assignee.email
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return nil
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			log.Printf("Row scan failed: %v", err)
			continue
		}
		emails = append(emails, email)
	}

	return emails
}
