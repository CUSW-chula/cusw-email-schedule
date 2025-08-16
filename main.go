package main

import (
	"database/sql"
	"log"
	"task-scheduler/lib"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to the database
	db := lib.ConnectDB()
	defer db.Close()

	// Send initial test email on app startup
	log.Println("📩 Sending initial test emails...")

	// Always send test email on startup
	testTasks := []lib.Task{
		{
			ID:            "cmb6blo4y017gpm1q50ku5vdx",
			Title:         "ระบบแจ้งเตือนงานเริ่มทำงานแล้ว",
			Description:   "ทดสอบระบบส่งอีเมลแจ้งเตือนงาน - คลิกปุ่ม View Task เพื่อไปยังหน้างาน",
			Status:        "Assigned",
			ProjectID:     "test-project-1",
			ProjectTitle:  "CUSW Email Scheduler",
			AssigneeName:  "Test User",
			AssignorName:  "System Administrator",
			AssigneeEmail: "bunyawatapp37204@gmail.com",
		},
	}

	// Send test email
	if err := lib.SendEmail("bunyawatapp37204@gmail.com", testTasks); err != nil {
		log.Printf("❌ Failed to send test email: %v", err)
	} else {
		log.Printf("✅ Test email sent successfully")
	}

	// Set up Cron Scheduler
	c := cron.New()

	// Define job to send emails for tasks due in 3 days (runs at 9:00 AM)
	_, err := c.AddFunc("0 9 * * *", func() {
		log.Println("⏰ Cron Job Started: Sending 3-day reminders...")
		sendRemindersForDays(db, 3)
	})
	if err != nil {
		log.Fatalf("❌ Failed to set up 3-day reminder Cron Job: %v", err)
	}

	// Define job to send emails for tasks due in 2 days (runs at 2:00 PM)
	_, err = c.AddFunc("0 14 * * *", func() {
		log.Println("⏰ Cron Job Started: Sending 2-day reminders...")
		sendRemindersForDays(db, 2)
	})
	if err != nil {
		log.Fatalf("❌ Failed to set up 2-day reminder Cron Job: %v", err)
	}

	// Define job to send emails for tasks due in 1 day (runs at 4:00 PM)
	_, err = c.AddFunc("0 16 * * *", func() {
		log.Println("⏰ Cron Job Started: Sending 1-day reminders...")
		sendRemindersForDays(db, 1)
	})
	if err != nil {
		log.Fatalf("❌ Failed to set up 1-day reminder Cron Job: %v", err)
	}

	// Start the Cron scheduler
	c.Start()
	log.Println("✅ Email notification system is running...")
	log.Println("📅 Schedule:")
	log.Println("   - 3-day reminders at 09:00")
	log.Println("   - 2-day reminders at 14:00")
	log.Println("   - 1-day reminders at 16:00")

	// Prevent the program from exiting
	select {}
}

// sendRemindersForDays ส่งอีเมลแจ้งเตือนสำหรับงานที่จะครบกำหนดในจำนวนวันที่กำหนด
func sendRemindersForDays(db *sql.DB, days int) {
	// Get all users with tasks due in specified days
	userEmails := lib.GetUniqueUserEmailsDueInDays(db, days)
	log.Printf("📋 Found %d users with tasks due in %d day(s)", len(userEmails), days)

	// Send email to each user
	for _, email := range userEmails {
		tasks := lib.QueryTasksByUserEmailDueInDays(db, email, days)
		if len(tasks) > 0 {
			log.Printf("📨 Sending %d-day reminder to: %s (%d tasks)", days, email, len(tasks))
			if err := lib.SendEmail(email, tasks); err != nil {
				log.Printf("❌ Failed to send %d-day reminder to %s: %v", days, email, err)
			} else {
				log.Printf("✅ %d-day reminder sent successfully to: %s", days, email)
			}
		}
	}
}
