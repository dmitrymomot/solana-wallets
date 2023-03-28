package main

import (
	"database/sql"

	"github.com/dmitrymomot/oauth2-server/lib/client"
	"github.com/dmitrymomot/oauth2-server/lib/middleware"
	"github.com/dmitrymomot/solana-wallets/internal/kitlog"
	"github.com/dmitrymomot/solana-wallets/internal/solanawallet"
	"github.com/dmitrymomot/solana-wallets/svc/balance"
	"github.com/dmitrymomot/solana-wallets/svc/wallet"
	wallet_repository "github.com/dmitrymomot/solana-wallets/svc/wallet/repository"
	solanaClient "github.com/dmitrymomot/solana/client"
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

	// Init HTTP router
	r := initRouter(logger.WithField("component", "http-router"))

	// OAuth2 middleware
	oauth2Mdw := middleware.GokitAuthMiddleware(client.Introspect(oauth2IntrospectURL))

	// Init Solana client
	solClient := solanaClient.New(solanaClient.SetSolanaEndpoint(solanaRPCURL))

	// Init wallet service
	{
		repo, err := wallet_repository.Prepare(ctx, db)
		if err != nil {
			logger.WithError(err).Fatal("Failed to prepare wallet repository")
		}
		r.Mount("/wallet", wallet.MakeHTTPHandler(
			wallet.MakeEndpoints(wallet.NewService(
				repo,
				solanawallet.NewClient(walletSecretSalt),
				solClient,
			), oauth2Mdw),
			kitlog.NewLogger(logger.WithField("component", "wallet-service")),
		))
	}

	// Init balance service
	{
		r.Mount("/balance", balance.MakeHTTPHandler(
			balance.MakeEndpoints(
				balance.NewService(solClient),
				oauth2Mdw,
			),
			kitlog.NewLogger(logger.WithField("component", "balance-service")),
		))
	}

	// Run HTTP server
	eg.Go(runServer(ctx, httpPort, r, logger.WithField("component", "http-server")))

	// Run all goroutines
	if err := eg.Wait(); err != nil {
		logger.WithError(err).Fatal("Error occurred")
	}
}
