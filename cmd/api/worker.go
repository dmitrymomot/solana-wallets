package main

import (
	"github.com/hibiken/asynq"
)

type (
	// taskHandler is an interface for task handlers.
	taskHandler interface {
		Register(*asynq.ServeMux)
	}
)

// setupQueue creates a new queue client and registers task handlers.
func runQueueServer(redisConnOpt asynq.RedisConnOpt, log asynq.Logger, handlers ...taskHandler) func() error {
	return func() error {
		// Setup asynq server
		srv := asynq.NewServer(
			redisConnOpt,
			asynq.Config{
				Concurrency:     workerConcurrency,
				Logger:          log,
				LogLevel:        getAsynqLogLevel(workerLogLevel),
				ShutdownTimeout: workerShutdownTimeout,
				Queues: map[string]int{
					queueName: workerConcurrency,
				},
			},
		)

		// Run server
		return srv.Run(registerQueueHandlers(handlers...))
	}
}

// registerQueueHandlers registers handlers for each task type.
func registerQueueHandlers(handlers ...taskHandler) *asynq.ServeMux {
	mux := asynq.NewServeMux()

	// Register handlers
	for _, h := range handlers {
		h.Register(mux)
	}

	return mux
}

type (
	schedulerHandler interface {
		Schedule(*asynq.Scheduler)
	}
)

// runScheduler creates a new scheduler server and registers task handlers.
func runScheduler(redisConnOpt asynq.RedisConnOpt, log asynq.Logger, handlers ...schedulerHandler) func() error {
	return func() error {
		// Setup asynq scheduler
		scheduler := asynq.NewScheduler(
			redisConnOpt,
			&asynq.SchedulerOpts{
				Logger:   log,
				LogLevel: getAsynqLogLevel(workerLogLevel),
			},
		)

		// Register handlers
		for _, h := range handlers {
			h.Schedule(scheduler)
		}

		// Run scheduler
		return scheduler.Run()
	}
}

// get asynq log level by string.
func getAsynqLogLevel(level string) asynq.LogLevel {
	switch level {
	case "debug":
		return asynq.DebugLevel
	case "info":
		return asynq.InfoLevel
	case "warn":
		return asynq.WarnLevel
	case "error":
		return asynq.ErrorLevel
	case "fatal":
		return asynq.FatalLevel
	default:
		return asynq.InfoLevel
	}
}
