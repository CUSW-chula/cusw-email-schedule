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

	// Set up Cron Scheduler
	c := cron.New()

	// Define job to send email every day at 8:00 AM
	_, err := c.AddFunc("0 8 * * *", func() {
		log.Println("Starting to fetch tasks...")
		tasks := lib.QueryTasks(db)
		log.Printf("Found %d tasks to notify\n", len(tasks))

		// Group tasks by email
		tasksByEmail := make(map[string][]lib.Task)
		for _, task := range tasks {
			tasksByEmail[task.Email] = append(tasksByEmail[task.Email], task)
		}

		// Send email to each user
		for email, userTasks := range tasksByEmail {
			log.Printf("Sending email to: %s (%d tasks)", email, len(userTasks))
			if err := lib.SendEmail(email, userTasks); err != nil {
				log.Printf("Failed to send email: %v", err)
			} else {
				log.Printf("✅ Email sent successfully to: %s", email)
			}
		}
	})

	if err != nil {
		log.Fatalf("Failed to set up Cron Job: %v", err)
	}

	c.Start()
	log.Println("Email notification system is running... Waiting for the next execution at 08:00.")

	// Wait for a signal to prevent the program from exiting
	select {}
}
