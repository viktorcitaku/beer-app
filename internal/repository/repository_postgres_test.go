package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/viktorcitaku/beer-app/internal/db"
)

type postgresContainer struct {
	testcontainers.Container
}

func setupPostgres(ctx context.Context, dbname, user, password string) (*postgresContainer, error) {
	abs, err := filepath.Abs("./../../../../scripts/db/")
	if err != nil {
		return nil, err
	}

	logStrategy := wait.NewLogStrategy("database system is ready to accept connections")
	logStrategy.WithOccurrence(2)

	req := testcontainers.ContainerRequest{
		Image: "postgres:alpine",
		Env: map[string]string{
			"POSTGRES_DB":       dbname,
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": password,
		},
		Mounts: []testcontainers.ContainerMount{
			{
				Source: testcontainers.GenericBindMountSource{
					HostPath: abs,
				},
				Target:   "/docker-entrypoint-initdb.d",
				ReadOnly: false,
			},
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForAll(logStrategy).WithDeadline(60 * time.Second),
	}

	var container testcontainers.Container
	container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &postgresContainer{container}, nil
}

func TestReadFromDatabase(t *testing.T) {
	ctx := context.Background()
	const dbname = "beer"
	const user = "test"
	const password = "test"
	postgres, err := setupPostgres(ctx, dbname, user, password)
	t.Cleanup(func() {
		if err := postgres.Terminate(ctx); err != nil {
			t.Error(err)
		}
	})
	if err != nil {
		t.Error(err)
	}

	host, err := postgres.Host(ctx)
	if err != nil {
		t.Error(err)
	}

	port, err := postgres.MappedPort(ctx, "5432")
	if err != nil {
		t.Error(err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port.Port(), user, password, dbname)

	conn := db.ConnectDatabase(connStr)
	defer db.CloseConnection()

	rows, err := conn.Query("SELECT * FROM beer.user_profile")
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			t.Error(err)
		}
	}(rows)
	if err != nil {
		t.Error(err)
	}

	for rows.Next() {
		var (
			email      string
			lastUpdate time.Time
		)

		if err := rows.Scan(&email, &lastUpdate); err != nil {
			t.Error(err)
		}

		fmt.Printf("User: %v | %v\n", email, lastUpdate)
	}
}
