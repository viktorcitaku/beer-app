package db_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/viktorcitaku/beer-app/internal/db"
	"go.mongodb.org/mongo-driver/bson"
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

func TestReadWriteMongo(t *testing.T) {
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

	coll := client.Database("school").Collection("students")
	address1 := Address{"1 Lakewood Way", "Elwood City", "PA"}
	student1 := Student{FirstName: "Arthur", Address: address1, Age: 8}
	res, err := coll.InsertOne(ctx, student1)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Mongo Client InsertOne result: %v\n", res.InsertedID)

	var result bson.D
	if err = coll.FindOne(ctx, bson.D{{"first_name", "Arthur"}}).Decode(&result); err != nil {
		t.Error(err)
	}

	fmt.Printf("Mongo Client FindOne result: %v\n", result)
}

type Address struct {
	Street string
	City   string
	State  string
}
type Student struct {
	FirstName string  `bson:"first_name,omitempty"`
	LastName  string  `bson:"last_name,omitempty"`
	Address   Address `bson:"inline"`
	Age       int
}
