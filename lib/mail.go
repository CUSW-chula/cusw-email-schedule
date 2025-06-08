package lib

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"

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
	subject := "🔔 แจ้งเตือนงานใกล้ถึงกำหนด"
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
			body { font-family: Arial, sans-serif; }
			.task-list { margin: 20px 0; }
			.task-item { padding: 10px; border-bottom: 1px solid #eee; }
			.deadline { color: #d9534f; font-weight: bold; }
		</style>
	</head>
	<body>
		<h2>📋 งานที่กำลังจะถึงกำหนด</h2>
		<div class="task-list">
			{{range .}}
			<div class="task-item">
				<p><strong>{{.Task}}</strong></p>
				<p class="deadline">⏰ กำหนดส่ง: {{.Deadline.Format "2006-01-02 15:04"}}</p>
			</div>
			{{end}}
		</div>
		<p>ด้วยความนับถือ,<br>ทีมงาน Task Scheduler</p>
	</body>
	</html>
	`
	t, _ := template.New("email").Parse(tmpl)
	var buf bytes.Buffer
	t.Execute(&buf, tasks)
	return buf.String()
}