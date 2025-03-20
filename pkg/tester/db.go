package tester

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/db"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun"
)

const (
	postgresImage           = "postgres:16"
	postgresUser            = "mountain"
	postgresPassword        = "pass123"
	postgresDefaultDatabase = "glacier"
)

type Database string

var dbOnce sync.Once
var dbInstance *bun.DB

func DB() *bun.DB {
	dbOnce.Do(initDB)

	return dbInstance
}

// initDB a postgres database via testcontainers
func initDB() {
	c := context.Background()

	host := "localhost"

	url := func(host string, port nat.Port) string {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPassword, host, port.Port(), postgresDefaultDatabase)
	}

	req := testcontainers.ContainerRequest{
		Image:        postgresImage,
		Hostname:     host,
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForSQL("5432/tcp", "postgres", url).WithStartupTimeout(1 * time.Minute),
		Env: map[string]string{
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_PASSWORD": postgresPassword,
			"POSTGRES_DB":       postgresDefaultDatabase,
		},
	}

	postgres, err := testcontainers.GenericContainer(c, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatalf("can't create postgres postgres: %v", err)
	}

	port, err := postgres.MappedPort(c, "5432/tcp")

	if err != nil {
		log.Fatalf("can't get mapped port postgres: %v", err)
	}

	log.Printf("started postgres container port=%s user=%s pass=%s", port.Port(), postgresUser, postgresPassword)

	dbInstance = mustMigrate(host, port.Port())
}

// mustMigrate runs a migrations and returns a connection
func mustMigrate(host string, port string) *bun.DB {
	// find the current directory so we can infer the migrations folder
	_, filename, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(filename)
	if err != nil {
		log.Fatalf("failed to get absolute path: %v", err)
	}

	cmdName := "bash"
	migrationScriptPath := filepath.Join(dir, "../../../scripts/migrate.sh")

	// on windows, we must execute the script with git bash
	if runtime.GOOS == "windows" {
		cmdName = "C:\\Program Files\\Git\\bin\\bash.exe"
	}

	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port)
	os.Setenv("DB_NAME", postgresDefaultDatabase)
	os.Setenv("DB_USER", postgresUser)
	os.Setenv("DB_PASS", postgresPassword)

	cmd := exec.Command(cmdName, migrationScriptPath, "--db-host", host, "--db-port", port)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatalf("failed to create stdout pipe: %v", err)
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()

	if err != nil {
		log.Fatalf("failed to create stderr pipe: %v", err)
	}
	defer stderr.Close()

	err = cmd.Start()

	if err != nil {
		log.Fatalf("failed to start command: %v", err)
	}

	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for scanner.Scan() {
		log.Println(scanner.Text())
	}

	err = cmd.Wait()

	if err != nil {
		log.Fatalf("command failed: %v", err)
	}

	return mustConnectDB(host, port, postgresDefaultDatabase)
}

// mustConnectDB creates a db connection or panics
func mustConnectDB(host, port, database string) *bun.DB {
	d, err := db.Postgres(host, port, postgresUser, postgresPassword, database)

	if err != nil {
		log.Fatalf("unable to connect to db: %v", err)
	}

	return d
}
