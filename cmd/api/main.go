package main

import (
	"database/sql"

	"github.com/hibiken/asynq"
	_ "github.com/lib/pq" // init pg driver
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Init logger with fields
	logger := logrus.WithFields(logrus.Fields{
		"app":       appName,
		"build_tag": buildTagRuntime,
	})
	defer func() { logger.Info("Server successfully shutdown") }()

	// Errgroup with context
	eg, ctx := errgroup.WithContext(newCtx(logger.WithField("component", "context")))

	// Init DB connection
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		logger.WithError(err).Fatal("Failed to init db connection")
	}
	defer db.Close()

	db.SetMaxOpenConns(dbMaxOpenConns)
	db.SetMaxIdleConns(dbMaxIdleConns)

	if err := db.Ping(); err != nil {
		logger.WithError(err).Fatal("Failed to ping db")
	}

	// Redis connect options for asynq client
	redisConnOpt, err := asynq.ParseRedisURI(redisConnString)
	if err != nil {
		logger.WithError(err).Fatal("failed to parse redis connection string")
	}

	// Init asynq client
	asynqClient := asynq.NewClient(redisConnOpt)
	defer asynqClient.Close()

	// Run asynq worker
	eg.Go(runQueueServer(
		redisConnOpt,
		logger.WithField("component", "queue-worker"),
		// TODO: add workers
	))

	// Run asynq scheduler
	eg.Go(runScheduler(
		redisConnOpt,
		logger.WithField("component", "scheduler"),
		// TODO: add schedulers
	))

	// Init HTTP router
	r := initRouter(logger.WithField("component", "http-router"))

	// TODO: add routes

	// Run HTTP server
	eg.Go(runServer(ctx, httpPort, r, logger.WithField("component", "http-server")))

	// Run all goroutines
	if err := eg.Wait(); err != nil {
		logger.WithError(err).Fatal("Error occurred")
	}
}
