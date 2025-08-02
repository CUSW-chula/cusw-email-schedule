package main

import (
	"log"
	"task-scheduler/lib"
	"time"

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
	testTasks := []lib.Task{
		{
			Email:    "bunyawatapp37204@gmail.com",
			Title:    "ระบบแจ้งเตือนงานเริ่มทำงานแล้ว",
			Deadline: time.Now().Add(24 * time.Hour),
		},
		{
			Email:    "melodymui2003@gmail.com",
			Title:    "ระบบแจ้งเตือนงานเริ่มทำงานแล้ว",
			Deadline: time.Now().Add(24 * time.Hour),
		},
		{
			Email:    "pond.phongsakorn1654@gmail.com",
			Title:    "ระบบแจ้งเตือนงานเริ่มทำงานแล้ว",
			Deadline: time.Now().Add(24 * time.Hour),
		},
	}

	// Group tasks by email
	tasksByEmail := make(map[string][]lib.Task)
	for _, task := range testTasks {
		tasksByEmail[task.Email] = append(tasksByEmail[task.Email], task)
	}

	// Send email to each recipient
	for email, tasks := range tasksByEmail {
		if err := lib.SendEmail(email, tasks); err != nil {
			log.Printf("❌ Failed to send initial email to %s: %v", email, err)
		} else {
			log.Printf("✅ Initial email sent successfully to: %s", email)
		}
	}

	// Set up Cron Scheduler
	c := cron.New()

	// Define job to send email every day at 4:00 PM
	_, err := c.AddFunc("0 16 * * *", func() {
		log.Println("⏰ Cron Job Started: Fetching tasks...")

		tasks := lib.QueryTasks(db)
		log.Printf("📋 Found %d tasks to notify", len(tasks))

		// Group tasks by email
		tasksByEmail := make(map[string][]lib.Task)
		for _, task := range tasks {
			tasksByEmail[task.Email] = append(tasksByEmail[task.Email], task)
		}

		// Send email to each user
		for email, userTasks := range tasksByEmail {
			log.Printf("📨 Sending email to: %s (%d tasks)", email, len(userTasks))
			if err := lib.SendEmail(email, userTasks); err != nil {
				log.Printf("❌ Failed to send email to %s: %v", email, err)
			} else {
				log.Printf("✅ Email sent successfully to: %s", email)
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
