package scheduler

import (
	"context"
	"encoding/json"
	"net/http"

	"sapasora/platform/logger"
)

// DashboardHandler provides HTTP handlers for the job dashboard
type DashboardHandler struct {
	queue      JobQueue
	workerPool WorkerPool
	log        *logger.Logger
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(
	queue JobQueue,
	workerPool WorkerPool,
	log *logger.Logger,
) *DashboardHandler {
	return &DashboardHandler{
		queue:      queue,
		workerPool: workerPool,
		log:        log,
	}
}

// RegisterRoutes registers the dashboard routes with the provided router
func (d *DashboardHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/admin/jobs", d.jobsPage)
	router.HandleFunc("/admin/jobs/stats", d.statsAPI)
	router.HandleFunc("/admin/jobs/retry", d.retryJob)
	router.HandleFunc("/admin/jobs/delete", d.deleteJob)
	router.HandleFunc("/admin/jobs/queues", d.queuesPage)
	router.HandleFunc("/admin/jobs/workers", d.workersPage)
}

// jobsPage renders the main jobs dashboard page
func (d *DashboardHandler) jobsPage(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, this would render a template page
	// For now, we'll return a simple JSON response
	ctx := r.Context()

	// In a real implementation, we would get jobs from the storage
	// For now, we'll just return stats
	stats, err := d.queue.Stats(ctx)
	if err != nil {
		d.log.Error("Failed to get stats", "error", err)
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"stats":   stats,
		"message": "Jobs dashboard - in a full implementation, this would render a UI page",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		d.log.Error("Failed to encode jobs response", "error", err)
	}
}

// queuesPage renders the queues dashboard page
func (d *DashboardHandler) queuesPage(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, this would render a template page
	ctx := r.Context()

	stats, err := d.queue.Stats(ctx)
	if err != nil {
		d.log.Error("Failed to get queue stats", "error", err)
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"stats":   stats,
		"message": "Queues dashboard - in a full implementation, this would render a UI page",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		d.log.Error("Failed to encode queues response", "error", err)
	}
}

// workersPage renders the workers dashboard page
func (d *DashboardHandler) workersPage(w http.ResponseWriter, r *http.Request) {
	workerStats := d.workerPool.Stats()

	response := map[string]any{
		"workerStats": workerStats,
		"message":     "Workers dashboard - in a full implementation, this would render a UI page",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		d.log.Error("Failed to encode workers response", "error", err)
	}
}

// statsAPI returns JSON statistics about the job system
func (d *DashboardHandler) statsAPI(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats, err := d.queue.Stats(ctx)
	if err != nil {
		d.log.Error("Failed to get stats", "error", err)
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	workerStats := d.workerPool.Stats()

	response := map[string]any{
		"queue_stats":  stats,
		"worker_stats": workerStats,
		"timestamp": r.Context().
			Value("timestamp"),
		// In a real implementation, this would be set
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		d.log.Error("Failed to encode stats response", "error", err)
	}
}

// retryJob handles retrying a failed job
func (d *DashboardHandler) retryJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobID := r.URL.Query().Get("id")
	if jobID == "" {
		http.Error(w, "Job ID required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := d.queue.Retry(ctx, jobID); err != nil {
		d.log.Error("Failed to retry job", "job_id", jobID, "error", err)
		http.Error(w, "Failed to retry job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Job retried successfully",
	})
}

// deleteJob handles deleting a job
func (d *DashboardHandler) deleteJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobID := r.URL.Query().Get("id")
	if jobID == "" {
		http.Error(w, "Job ID required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := d.queue.Delete(ctx, jobID); err != nil {
		d.log.Error("Failed to delete job", "job_id", jobID, "error", err)
		http.Error(w, "Failed to delete job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Job deleted successfully",
	})
}

// getAllJobs gets all jobs from storage (in a real implementation, this might be paginated)
func (d *DashboardHandler) getAllJobs(ctx context.Context) ([]*Job, error) {
	// In a real implementation, we'd need access to the storage layer
	// For now, we'll return an empty slice to avoid compilation errors
	// The actual implementation would depend on how we access the storage
	return []*Job{}, nil
}

// Inertia views would be implemented as templ files in the resources/views directory
// For now, we'll just define the necessary functions
