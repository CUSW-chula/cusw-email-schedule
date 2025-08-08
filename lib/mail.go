package lib

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
	"strings"

	"gopkg.in/gomail.v2"
)

// SendEmail ส่งอีเมลแจ้งเตือนงานให้ผู้ใช้
func SendEmail(email string, tasks []Task) error {
	// ดึงค่าการตั้งค่าจาก Environment Variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := 587
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	sender := os.Getenv("EMAIL_SENDER")

	// กำหนดเนื้อหาอีเมลแบบ HTML
	subject := "🔔 Task Notification: CUSW Workspace"
	if len(tasks) > 0 {
		subject = fmt.Sprintf("🔔 Task Notification: %s – %s", tasks[0].Title, tasks[0].ProjectTitle)
	}
	body := buildEmailBody(tasks)

	// สร้างอีเมล
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// ตั้งค่า SMTP Dialer
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// ส่งอีเมล
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("ส่งอีเมลไม่สำเร็จ: %v", err)
	}
	return nil
}

func buildEmailBody(tasks []Task) string {
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body { 
				font-family: Arial, sans-serif; 
				line-height: 1.6;
				color: #333;
				max-width: 600px;
				margin: 0 auto;
				padding: 20px;
			}
			.header {
				background-color: #f8f9fa;
				padding: 20px;
				border-radius: 8px;
				margin-bottom: 20px;
				text-align: center;
			}
			.logo {
				max-width: 200px;
				height: auto;
				margin-bottom: 15px;
				display: block;
				margin-left: auto;
				margin-right: auto;
				/* Ensure SVG displays properly in email clients */
				width: 200px;
				max-height: 100px;
			}
			.header h2 {
				margin: 0;
				color: #2c3e50;
				font-size: 20px;
			}
			.task-container {
				background-color: #ffffff;
				border: 1px solid #e9ecef;
				border-radius: 8px;
				padding: 20px;
				margin-bottom: 20px;
			}
			.task-info {
				margin: 10px 0;
			}
			.task-info strong {
				color: #495057;
			}
			.description-content {
				margin-top: 8px;
				line-height: 1.6;
			}
			.description-content p {
				margin: 8px 0;
			}
			.status {
				display: inline-block;
				padding: 4px 12px;
				border-radius: 20px;
				font-size: 12px;
				font-weight: bold;
				text-transform: uppercase;
			}
			.status.assigned { background-color: #d1ecf1; color: #0c5460; }
			.status.underreview { background-color: #fff3cd; color: #856404; }
			.status.inrecheck { background-color: #ffeaa7; color: #b8860b; }
			.status.done { background-color: #d4edda; color: #155724; }
			.status.unassigned { background-color: #f8d7da; color: #721c24; }
			.button {
				display: inline-block;
				background-color: #007bff;
				color: white !important;
				padding: 12px 24px;
				text-decoration: none;
				border-radius: 5px;
				margin: 20px 0;
				font-weight: bold;
				border: none;
				cursor: pointer;
				transition: background-color 0.3s;
			}
			.button:hover {
				background-color: #0056b3;
				color: white !important;
				text-decoration: none;
			}
			.footer {
				margin-top: 30px;
				padding-top: 20px;
				border-top: 1px solid #e9ecef;
				font-size: 12px;
				color: #6c757d;
				text-align: center;
			}
			.footer-logo {
				max-width: 60px;
				height: auto;
				margin-bottom: 10px;
				opacity: 0.7;
				display: block;
				margin-left: auto;
				margin-right: auto;
				/* Ensure SVG displays properly in email clients */
				width: 60px;
				max-height: 30px;
			}
		</style>
	</head>
	<body>
		{{range .}}
		<div class="header">
			<img src="https://cusw-workspace.sa.chula.ac.th/asset/logo/s2.png" alt="CUSW Logo" class="logo" style="width: 200px; max-height: 100px; display: block; margin: 0 auto 15px;">
			<h2>🔔 Task Notification: {{.Title}} – {{.ProjectTitle}}</h2>
		</div>
		
		<p>Hi {{.AssigneeName}},</p>
		
		<p>You have a new update regarding a task in the cusw-workspace platform:</p>
		
		<div class="task-container">
			<div class="task-info">
				<strong>📌 Task:</strong> {{.Title}}
			</div>
			<div class="task-info">
				<strong>📁 Project:</strong> {{.ProjectTitle}}
			</div>
			<div class="task-info">
				<strong>🗓 Due Date:</strong> {{if .EndDate}}{{.EndDate.Format "January 2, 2006 at 15:04"}}{{else}}Not specified{{end}}
			</div>
			<div class="task-info">
				<strong>👤 Assigned By:</strong> {{.AssignorName}}
			</div>
			<div class="task-info">
				<strong>📎 Status:</strong> <span class="status {{.Status | lower}}">{{.Status}}</span>
			</div>
			
			{{if .Description}}
			<div class="task-info">
				<strong>📝 Description:</strong><br>
				<div class="description-content">{{.Description | safeHTML}}</div>
			</div>
			{{end}}
		</div>
		
		<p>You can view the task and take necessary actions by clicking the button below:</p>
		
		<a href="https://cusw-workspace.sa.chula.ac.th/tasks/{{.ID}}" class="button" target="_blank">👉 View Task</a>
		{{end}}
		
		<div class="footer">
			<img src="https://cusw-workspace.sa.chula.ac.th/asset/logo/s2.png" alt="CUSW Logo" class="footer-logo" style="width: 60px; max-height: 30px; display: block; margin: 0 auto 10px; opacity: 0.7;">
			<p>Thank you,<br>
			<strong>CUSW+</strong><br>
			<em>This is an automated message from your task workspace.</em></p>
		</div>
	</body>
	</html>
	`

	// Create a function map for template helpers
	funcMap := template.FuncMap{
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	t, _ := template.New("email").Funcs(funcMap).Parse(tmpl)
	var buf bytes.Buffer
	t.Execute(&buf, tasks)
	return buf.String()
}

// SendTaskNotification ส่งอีเมลแจ้งเตือนสำหรับ task เดียว
func SendTaskNotification(task Task) error {
	return SendEmail(task.AssigneeEmail, []Task{task})
}
