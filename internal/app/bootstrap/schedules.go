package bootstrap

import (
	"sapasora/platform/logger"
	"sapasora/platform/scheduler"
)

// NewSchedules creates a new cron scheduler
// @wired:provide(group=schedules)
func NewSchedules(log *logger.Logger) scheduler.Schedules {
	return func(cron *scheduler.CRON) {
		/* // notice that we have 6 parts instead of 5
		 * // that's because we use seconds instead of minutes
		 * // which give us more precises control over the schedule
		 * err := cron.Add("0 * * * * *", "ping", func() error {
		 * 	log.Info("pong")
		 * 	return nil
		 * })
		 * if err != nil {
		 * 	log.Error("Failed to schedule job", "error", err)
		 * }
		 */
	}
}
