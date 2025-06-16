package main

import (
	"context"
	"generate_technical_challenge_2025/internal/database"
	"generate_technical_challenge_2025/internal/handler"
	"generate_technical_challenge_2025/internal/server"
	"generate_technical_challenge_2025/internal/services"
	"generate_technical_challenge_2025/internal/transactions"
	"generate_technical_challenge_2025/internal/utils"
	"log/slog"

	"github.com/sethvargo/go-envconfig"
)

func main() {
	// Setup logger and environment variables.
	logger := slog.New(slog.Default().Handler())

	logger.Info("Loading environment variables...")
	env := loadEnv()

	logger.Info("Creating database from environment variables...")
	db := database.CreateDatabase(env, logger)

	logger.Info("Auto migrating database schemas...")
	database.AutoMigrate(db)

	logger.Info("Initializing transaction layer...")
	transactions := transactions.CreateTransactions(logger, db)

	logger.Info("Intializing service layer...")
	services := services.CreateServices(logger, transactions)

	logger.Info("Intializing handler layer...")
	h := handler.CreateHandler(logger, services)

	port := "8081"

	logger.Info("Attaching handler and running server on http://localhost:" + port + "...")
	server.RunServer(h, ":"+port)
}

// Loads the environment variables as an EnvConfig
func loadEnv() utils.EnvConfig {
	var config utils.EnvConfig
	envFun := func() error { return envconfig.Process(context.Background(), &config) }
	utils.SafeCallErrorSupplier(envFun)
	return config
}
