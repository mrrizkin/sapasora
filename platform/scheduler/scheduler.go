package scheduler

import (
	"context"
	"time"

	"sapasora/platform/config"
	"sapasora/platform/database"
	"sapasora/platform/logger"

	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

// CRON handles scheduling and execution of cron jobs (traditional cron functionality)
type CRON struct {
	cron *cron.Cron
	log  *logger.Logger
}

// SchedulerIn contains the dependencies for the traditional Scheduler
type SchedulerIn struct {
	fx.In
	Logger    *logger.Logger
	Schedules []Schedules `group:"schedules"`
}

type Schedules func(c *CRON)

// NewScheduler creates a new traditional Scheduler instance and registers it in the fx lifecycle
func NewScheduler(in SchedulerIn) *CRON {
	s := &CRON{
		cron: cron.New(cron.WithSeconds()),
		log:  in.Logger,
	}

	for _, schedule := range in.Schedules {
		schedule(s)
	}

	return s
}

// Add schedules a new task with the given spec and name
func (s *CRON) Add(spec, name string, task func() error) error {
	_, err := s.cron.AddFunc(spec, func() {
		s.log.Info("Starting job execution", "job", name)
		startTime := time.Now()

		if err := task(); err != nil {
			s.log.Error(
				"Job execution failed",
				"job",
				name,
				"error",
				err,
				"duration",
				time.Since(startTime),
			)
		} else {
			s.log.Info("Job execution succeeded", "job", name, "duration", time.Since(startTime))
		}
	})
	return err
}

// Start begins the scheduler and logs the number of entries
func (s *CRON) Start() error {
	entryCount := len(s.cron.Entries())
	s.log.Info("Starting scheduler", "jobs", entryCount)

	if entryCount == 0 {
		s.log.Warn("Scheduler started with no jobs")
	}

	s.cron.Start()
	return nil
}

// Stop halts the scheduler and waits for jobs to complete
func (s *CRON) Stop() error {
	s.log.Info("Stopping scheduler", "jobs", len(s.cron.Entries()))
	ctx := s.cron.Stop()
	<-ctx.Done()
	return nil
}

// NewAdvancedSchedulerFromLifecycle creates a new AdvancedScheduler with lifecycle management
func NewAdvancedSchedulerFromLifecycle(
	cfg config.Config,
	log *logger.Logger,
	db *database.Database,
	lc fx.Lifecycle,
) (*AdvancedScheduler, error) {
	scheduler, err := NewAdvancedScheduler(cfg, log, db)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return scheduler.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return scheduler.Stop(ctx)
		},
	})

	return scheduler, nil
}
