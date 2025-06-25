// Full integration tests for server.
package server_test

import (
	"context"
	"generate_technical_challenge_2025/internal/database"
	"generate_technical_challenge_2025/internal/handler"
	"generate_technical_challenge_2025/internal/server"
	"generate_technical_challenge_2025/internal/services"
	"generate_technical_challenge_2025/internal/transactions"
	"generate_technical_challenge_2025/internal/utils"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	PORT   = 8008
	LOGGER = slog.New(slog.Default().Handler())
	CLIENT = func() *utils.TestClient { return utils.CreateTestClient(PORT, LOGGER) }
)

func runServer() {
	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"
	dbPort := "5432"
	LOGGER.Info("Creating postgres container")
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}
	dbHostFn := func() (string, error) { return postgresContainer.Host(ctx) }
	dbHost := utils.SafeCall(dbHostFn)

	dbHostPortFn := func() (nat.Port, error) { return postgresContainer.MappedPort(ctx, nat.Port(dbPort)) }
	dbPort = utils.SafeCall(dbHostPortFn).Port()

	envConfig := &utils.EnvConfig{
		DB_HOST:     dbHost,
		DB_PORT:     dbPort,
		DB_USER:     dbUser,
		DB_PASSWORD: dbPassword,
		DB_NAME:     dbName,
		PORT:        PORT,
	}
	db := database.CreateDatabase(*envConfig, LOGGER)

	database.AutoMigrate(db)

	memberTransactions := transactions.CreateMemberTransactions(LOGGER, db)

	memberServices := services.CreateMemberService(LOGGER, memberTransactions)

	h := handler.CreateHandler(LOGGER, memberServices)
	server.RunServer(h, *envConfig, LOGGER)
}

func TestMain(m *testing.M) {
	LOGGER.Info("Starting test server in a seperate go routine..")
	go runServer()
	if !CLIENT().CheckServer(time.Second * 30) {
		os.Exit(1)
	}
	LOGGER.Info("Finished setting up mock postgres container and server...")
	LOGGER.Info("Running tests...")
	code := m.Run()
	os.Exit(code)
}

func TestHealthCheck(t *testing.T) {
	expectedBody := map[string]interface{}{
		"message": "OK",
	}

	testVerify := CLIENT().GET("/healthcheck")
	testVerify.AssertStatusCode(200, t).AssertBody(expectedBody, t)
}

func TestUserWithNonValidNUIDReceives400(t *testing.T) {
	client := CLIENT()
	client.AddBody(map[string]any{
		"email": "notavalidemail@gmail.com",
		"nuid":  "1231",
	})
	client.AddHeaders(map[string]string{
		"Content-Type": "application/json",
	})
	testVerify := client.POST("/api/v1/member/register")
	testVerify.AssertStatusCode(400, t).AssertBody(map[string]any{
		"message": "Not a valid northeastern email address.",
	}, t)
}

func TestUserWithBadNUIDReceives400(t *testing.T) {
	client := CLIENT()
	client.AddBody(map[string]any{
		"email": "somebody@northeastern.edu",
		"nuid":  "1231",
	})
	client.AddHeaders(map[string]string{
		"Content-Type": "application/json",
	})
	testVerify := client.POST("/api/v1/member/register")
	testVerify.AssertStatusCode(400, t).AssertBody(map[string]any{
		"message": "Not a valid NUID.",
	}, t)
}

func TestUserReceives201OnGoodRequest(t *testing.T) {
	client := CLIENT()
	client.AddBody(map[string]any{
		"email": "somebody@northeastern.edu",
		"nuid":  "123456789", // NUID is 9 characters long
	})
	client.AddHeaders(map[string]string{
		"Content-Type": "application/json",
	})

	testVerify := client.POST("/api/v1/member/register")
	pred := func(prop any) bool {
		s, ok := prop.(string)
		if !ok {
			return ok
		}
		_, err := uuid.Parse(s)
		return err == nil
	}
	testVerify.AssertStatusCode(201, t).AssertProperty("id", pred, t)
}
