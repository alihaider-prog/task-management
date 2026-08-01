package main

import (
	// "database/sql"
	// "fmt"
	// "encoding/hex"
	// "os"
	// "log"

	"task-management/internal/database"
	"task-management/internal/handler"
	"task-management/internal/middleware"
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

	// var jwtSecret []byte
	// secretHex := os.Getenv("JWT_SECRET")
	// if secretHex == "" {
	// 	log.Fatal("JWT_SECRET environment variable not set")
	// }
	// decoded, err := hex.DecodeString(secretHex)
	// if err != nil {
	// 	log.Fatal("invalid JWT_SECRET: must be hex-encoded")
	// }
	// if len(decoded) < 32 {
	// 	log.Fatal("JWT_SECRET must be at least 256 bits (32 bytes)")
	// }
	// jwtSecret = decoded

	// os.Setenv("JWT_SECRET", string(jwtSecret))

	//auth
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

	// tasks
	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)

	// workspaces
	workspaceRepo := repository.NewWorkspaceRepository(db)
	workspaceService := service.NewWorkspaceService(workspaceRepo)
	workspaceHandler := handler.NewWorkspaceHandler(workspaceService)

	// projects
	projectRepo := repository.NewProjectRepository(db)
	projectService := service.NewProjectService(projectRepo)
	projectHandler := handler.NewProjectHandler(projectService)

	// members
	memberRepo := repository.NewMemberRepository(db)
	memberService := service.NewMemberService(memberRepo)
	memberHandler := handler.NewMemberHandler(memberService)

	roleMiddleware := middleware.NewRoleMiddleware(workspaceRepo, projectRepo, taskRepo, memberRepo)

	// tasks (handler)
	taskHandler := handler.NewTaskHandler(taskService, roleMiddleware)

	r := router.SetupRouter(
		authHandler,
		taskHandler,
		workspaceHandler,
		projectHandler,
		memberHandler,
		roleMiddleware,
	)

	r.Run()
}
