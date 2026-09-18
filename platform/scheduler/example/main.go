package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"sapasora/internal/app"
	"sapasora/platform/scheduler"

	"go.uber.org/fx"
)

func main() {
	// Create an FX app with just the scheduler components
	app := fx.New(
		app.Module, // This includes all platform modules including scheduler
		fx.Invoke(useAdvancedScheduler),
	)

	// Start the app for a short time to demonstrate the scheduler
	startCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		log.Fatalf("Failed to start app: %v", err)
	}

	// Keep the app running briefly to see scheduler in action
	time.Sleep(3 * time.Second)

	// Stop the app
	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Printf("Error stopping app: %v", err)
	}
}

func useAdvancedScheduler(advancedScheduler *scheduler.AdvancedScheduler) {
	ctx := context.Background()

	// Example: Create a simple job function
	simpleJob := func(ctx context.Context, args ...interface{}) error {
		fmt.Printf("Executing simple job with args: %v at %v\n", args, time.Now())
		return nil
	}

	// Enqueue the job
	err := advancedScheduler.PerformAsync(ctx, simpleJob, "hello", "world", 123)
	if err != nil {
		fmt.Printf("Error enqueuing job: %v\n", err)
	} else {
		fmt.Println("Successfully enqueued simple job")
	}

	// Example: Enqueue a job with delay using a goroutine approach (EnqueueIn is not fully implemented yet)
	delayedJob := func(ctx context.Context, args ...interface{}) error {
		fmt.Printf("Executing delayed job with args: %v at %v\n", args, time.Now())
		return nil
	}

	// For now, just enqueue immediately since the full delayed functionality requires more complex implementation
	err = advancedScheduler.PerformAsync(ctx, delayedJob, "delayed", "args")
	if err != nil {
		fmt.Printf("Error enqueuing job: %v\n", err)
	} else {
		fmt.Println("Successfully enqueued job")
	}

	// Display scheduler stats
	stats, err := advancedScheduler.GetStats(ctx)
	if err != nil {
		fmt.Printf("Error getting stats: %v\n", err)
	} else {
		fmt.Printf("Scheduler stats: %+v\n", stats)
	}

	workerStats := advancedScheduler.GetWorkerStats()
	fmt.Printf("Worker stats: %+v\n", workerStats)
}
