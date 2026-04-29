package tests

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"os"
	"path/filepath"
	"net"
	"google.golang.org/grpc"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

type TestSuite struct {
	PGContainer    *postgres.PostgresContainer
	RedisContainer *redis.RedisContainer
	DBURL          string
	RedisURL       string
	Ctx            context.Context
}

func SetupTestSuite() (*TestSuite, error) {
	ctx := context.Background()

	// Setup PostgreSQL
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	dbURL, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres connection string: %w", err)
	}

	// Setup Redis
	redisContainer, err := redis.Run(ctx, "redis:7-alpine")
	if err != nil {
		return nil, fmt.Errorf("failed to start redis container: %w", err)
	}

	redisURL, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get redis connection string: %w", err)
	}

	ts := &TestSuite{
		PGContainer:    pgContainer,
		RedisContainer: redisContainer,
		DBURL:          dbURL,
		RedisURL:       redisURL,
		Ctx:            ctx,
	}

	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get cwd: %w", err)
	}

	root := cwd
	for {
		if _, err := os.Stat(filepath.Join(root, "justfile")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			return nil, fmt.Errorf("could not find justfile in any parent directory")
		}
		root = parent
	}

	migrationsDir := filepath.Join(root, "pkg", "database", "migrations")
	if err := ts.RunMigrations(migrationsDir); err != nil {
		ts.Teardown()
		return nil, fmt.Errorf("failed to run migrations in %s: %w", migrationsDir, err)
	}

	return ts, nil
}

func (ts *TestSuite) RunMigrations(migrationsDir string) error {
	db, err := goose.OpenDBWithDriver("pgx", ts.DBURL)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer db.Close()

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}

func (ts *TestSuite) DBPool() *pgxpool.Pool {
	pool, _ := pgxpool.New(ts.Ctx, ts.DBURL)
	return pool
}

func (ts *TestSuite) RedisClient() *goredis.Client {
	opts, _ := goredis.ParseURL(ts.RedisURL)
	return goredis.NewClient(opts)
}

func (ts *TestSuite) Teardown() {
	if ts.PGContainer != nil {
		if err := ts.PGContainer.Terminate(ts.Ctx); err != nil {
			log.Printf("failed to terminate postgres container: %v", err)
		}
	}
	if ts.RedisContainer != nil {
		if err := ts.RedisContainer.Terminate(ts.Ctx); err != nil {
			log.Printf("failed to terminate redis container: %v", err)
		}
	}
}

type GRPCServerRunner interface {
	Serve(lis net.Listener) error
	Stop()
	GracefulStop()
}

func RunGRPCServer(server *grpc.Server) (string, error) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Printf("grpc server error: %v", err)
		}
	}()
	return lis.Addr().String(), nil
}
