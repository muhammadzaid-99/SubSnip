// internal/api/handlers.go
package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/muhammadzaid-99/SubSnip/internal/models"
	"github.com/muhammadzaid-99/SubSnip/internal/queue"
	"github.com/muhammadzaid-99/SubSnip/internal/status"
)

func SubmitTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Manual validation
	if task.VideoURL == "" || task.LanguageCode == "" || task.StartTime == "" || task.EndTime == "" || task.FramesToSkip <= 0 ||
		task.SubtitleBox.TopLeft.X == 0 && task.SubtitleBox.TopLeft.Y == 0 && task.SubtitleBox.BottomRight.X == 0 && task.SubtitleBox.BottomRight.Y == 0 {
		http.Error(w, "Missing or invalid required fields", http.StatusBadRequest)
		return
	}

	taskID := "task-" + uuid.NewString()
	status.Set(taskID, "queued")

	task.TaskID = taskID

	if err := queue.Publish(task); err != nil {
		status.Set(taskID, "failed")
		http.Error(w, "Failed to queue task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"task_id": taskID})
}

func TaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	st := status.Get(taskID)
	if st == "" {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": st})
}
