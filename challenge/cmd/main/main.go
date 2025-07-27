package main

import (
	"generate_technical_challenge_2025/internal/database"
	"generate_technical_challenge_2025/internal/handler"
	"generate_technical_challenge_2025/internal/server"
	"generate_technical_challenge_2025/internal/services"
	"generate_technical_challenge_2025/internal/transactions"
	"generate_technical_challenge_2025/internal/utils"
	"log/slog"
)

func main() {
	// Setup logger and environment variables.
	logger := slog.New(slog.Default().Handler())

	logger.Info("Loading environment variables...")
	env := utils.LoadEnv()

	logger.Info("Creating database from environment variables...")
	db := database.CreateDatabase(env, logger)

	logger.Info("Auto migrating database schemas...")
	database.AutoMigrate(db)

	logger.Info("Initializing transaction layer...")
	memberTransactions := transactions.CreateMemberTransactions(logger, db)
	challengeTransactions := transactions.CreateChallengeTransactions(logger, db)

	logger.Info("Intializing service layer...")
	memberServices := services.CreateMemberService(logger, memberTransactions)
	challengeServices := services.CreateChallengeService(
		logger, challengeTransactions)

	logger.Info("Intializing handler layer...")
	h := handler.CreateHandler(logger, memberServices, challengeServices)

	server.RunServer(h, env, logger, env.SLACK_WEBHOOK)
}
