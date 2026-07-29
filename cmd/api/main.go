package main

import (
	// "database/sql"
	// "fmt"
	// "log"

	"task-management/internal/database"
	"task-management/internal/handler"
	"task-management/internal/repository"
	"task-management/internal/router"
	"task-management/internal/service"

	// _ "github.com/lib/pq" // To register the driver.
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Initialize DB connection
	db := database.Init(".env")

	// // Use db.DB for database operations
	// rows, err := database.DB.Query("SELECT NOW()")
	// if err != nil {
	// 	log.Fatalf("Query failed: %v", err)
	// }
	// defer rows.Close()
	// var now string
	// for rows.Next() {
	// 	if err := rows.Scan(&now); err != nil {
	// 		log.Fatalf("Scan failed: %v", err)
	// 	}
	// 	fmt.Println("Current time from DB:", now)
	// }
	// if err := rows.Err(); err != nil {
	// 	log.Fatalf("rows failed: %v", err)
	// }
	// fmt.Println("connected")

	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	r := router.SetupRouter(taskHandler)

	r.Run()
}
