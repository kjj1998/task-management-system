package testutils

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MySQLContainer struct {
	Container testcontainers.Container
}

func CreateMySQLContainer(ctx context.Context) (*MySQLContainer, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting current directory: %v", err)
	}
	fmt.Println("Current working directory:", dir)
	absPath, err := filepath.Abs(filepath.Join("..", "testutils", "init-test-db.sql"))
	fmt.Println(absPath)
	if err != nil {
		log.Fatal(err)
	}
	r, err := os.Open(absPath)
	if err != nil {
		log.Fatal(err)
	}

	req := testcontainers.ContainerRequest{
		Image: "mysql:8.0",
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "password",
			"MYSQL_DATABASE":      "taskapi",
			"MYSQL_USER":          "testuser",
			"MYSQL_PASSWORD":      "testpass",
		},
		ExposedPorts: []string{"3306/tcp"},
		WaitingFor:   wait.ForLog("port: 3306  MySQL Community Server - GPL").WithStartupTimeout(30 * time.Second),
		Files: []testcontainers.ContainerFile{
			{
				Reader:            r,
				HostFilePath:      absPath,
				ContainerFilePath: "/docker-entrypoint-initdb.d/init-db.sql",
				FileMode:          0o755,
			},
		},
	}

	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &MySQLContainer{
		Container: mysqlC,
	}, nil
}
