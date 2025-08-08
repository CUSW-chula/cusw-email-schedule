package main

import (
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

	// Get users with upcoming tasks for testing
	userEmails := lib.GetUniqueUserEmails(db)
	if len(userEmails) == 0 {
		log.Println("ℹ️ No users with upcoming tasks found for testing")

		// Create test data if no real data exists
		testTasks := []lib.Task{
			{
				ID:            "test-1",
				Title:         "ระบบแจ้งเตือนงานเริ่มทำงานแล้ว",
				Description:   "ทดสอบระบบส่งอีเมลแจ้งเตือนงาน",
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
	} else {
		// Send real data
		for _, email := range userEmails {
			tasks := lib.QueryTasksByUserEmail(db, email)
			if len(tasks) > 0 {
				if err := lib.SendEmail(email, tasks); err != nil {
					log.Printf("❌ Failed to send initial email to %s: %v", email, err)
				} else {
					log.Printf("✅ Initial email sent successfully to: %s (%d tasks)", email, len(tasks))
				}
			}
		}
	}

	// Set up Cron Scheduler
	c := cron.New()

	// Define job to send email every day at 4:00 PM
	_, err := c.AddFunc("0 16 * * *", func() {
		log.Println("⏰ Cron Job Started: Fetching tasks...")

		// Get all users with upcoming tasks
		userEmails := lib.GetUniqueUserEmails(db)
		log.Printf("📋 Found %d users with upcoming tasks", len(userEmails))

		// Send email to each user
		for _, email := range userEmails {
			tasks := lib.QueryTasksByUserEmail(db, email)
			if len(tasks) > 0 {
				log.Printf("📨 Sending email to: %s (%d tasks)", email, len(tasks))
				if err := lib.SendEmail(email, tasks); err != nil {
					log.Printf("❌ Failed to send email to %s: %v", email, err)
				} else {
					log.Printf("✅ Email sent successfully to: %s", email)
				}
			}
		}
	})

	if err != nil {
		log.Fatalf("❌ Failed to set up Cron Job: %v", err)
	}

	// Start the Cron scheduler
	c.Start()
	log.Println("✅ Email notification system is running... Waiting for the next execution at 16:00.")

	// Prevent the program from exiting
	select {}
}
