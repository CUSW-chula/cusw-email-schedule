package main

import (
	"fmt"
	"log"
	"task-scheduler/lib"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get database connection
	db := lib.ConnectDB()
	defer db.Close()

	// Initialize cron scheduler
	c := cron.New()

	// Schedule daily job at 08:00 AM
	_, err := c.AddFunc("0 8 * * *", func() {
		log.Println("Executing scheduled query...")
		tasks := lib.QueryTasks(db)
		log.Printf("Retrieved %d tasks\n", len(tasks))
		
		// Process tasks (example: print to console)
		for i, task := range tasks {
			fmt.Printf("[%d] Email: %s, Deadline: %s, Task: %s\n",
				i+1,
				task.Email,
				task.Deadline.Format("2006-01-02"),
				task.Task)
		}
	})
	
	if err != nil {
		log.Fatalf("Error scheduling job: %v", err)
	}

	c.Start()
	log.Println("Scheduler started. Waiting for next run at 08:00 AM...")
	
	// Keep main thread alive
	select{}
}


