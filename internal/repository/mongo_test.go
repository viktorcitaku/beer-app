package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/viktorcitaku/beer-app/internal/db"
	"github.com/viktorcitaku/beer-app/internal/repository"
)

type mongoContainer struct {
	testcontainers.Container
}

func setupMongo(ctx context.Context) (*mongoContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Waiting for connections"),
			wait.ForListeningPort("27017/tcp"),
		).WithDeadline(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &mongoContainer{container}, nil
}

func TestMongoRepository(t *testing.T) {
	ctx := context.Background()
	mongo, err := setupMongo(ctx)
	t.Cleanup(func() {
		if err := mongo.Terminate(ctx); err != nil {
			t.Error(err)
		}
	})
	if err != nil {
		t.Error(err)
	}

	endpoint, err := mongo.Endpoint(ctx, "mongodb")
	if err != nil {
		t.Error(err)
	}

	client := db.ConnectMongo(ctx, endpoint)
	defer db.CloseMongoClient(ctx)

	if err = client.Ping(ctx, nil); err != nil {
		t.Error(err)
	}

	upr := repository.NewUserPreferencesRepository(ctx, client)

	up := &repository.UserPreferences{
		Id:                 1,
		BeerId:             1,
		BeerName:           "Wadler",
		UserEmail:          "test@test.com",
		DrunkTheBeerBefore: false,
		GotDrunk:           false,
		LastTime:           "2023-01-01",
		Rating:             0,
		Comment:            "No Comment",
	}
	err = upr.Save(up)
	if err != nil {
		t.Error(err)
	}

	up1, err := upr.FindById(1)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(up1)

	up2, err := upr.FindByBeerId(1)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(up2)

	ups, err := upr.FindByUserEmail("test@test.com")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(ups)
}
