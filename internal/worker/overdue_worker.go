package worker

import (
	"context"
	"log"
	"time"

	"task-management/internal/service"
)

type OverdueWorker struct {
	taskService *service.TaskService
	interval    time.Duration
}

func NewOverdueWorker(taskService *service.TaskService) *OverdueWorker {
	return &OverdueWorker{
		taskService: taskService,
		interval:    time.Minute,
	}
}

// Run starts the background worker loop. It returns when the provided context is cancelled.
// It processes overdue tasks in batches using `batchSize` and records `systemUserID` in history.
func (w *OverdueWorker) Run(ctx context.Context, batchSize int, systemUserID int64) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Run once immediately
	w.runOnce(batchSize, systemUserID)

	for {
		select {
		case <-ctx.Done():
			// Context cancelled — drain and exit
			log.Println("OverdueWorker: shutdown requested, exiting")
			return
		case <-ticker.C:
			w.runOnce(batchSize, systemUserID)
		}
	}
}

func (w *OverdueWorker) runOnce(batchSize int, systemUserID int64) {
	count, err := w.taskService.ProcessOverdue(batchSize, systemUserID)
	if err != nil {
		log.Printf("OverdueWorker: error processing overdue tasks: %v", err)
		return
	}
	if count > 0 {
		log.Printf("OverdueWorker: marked %d tasks overdue", count)
	}
}

