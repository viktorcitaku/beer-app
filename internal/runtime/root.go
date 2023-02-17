package runtime

import (
	"context"

	"github.com/viktorcitaku/beer-app/internal/beerapi"
	"github.com/viktorcitaku/beer-app/internal/cache"
	"github.com/viktorcitaku/beer-app/internal/config"
	controller "github.com/viktorcitaku/beer-app/internal/controller/v1"
	"github.com/viktorcitaku/beer-app/internal/db"
	"github.com/viktorcitaku/beer-app/internal/repository"
	"github.com/viktorcitaku/beer-app/internal/router"
	"github.com/viktorcitaku/beer-app/internal/server"
)

func Run() {
	ctx := context.Background()
	cfg := config.GetConfig()

	// Initialize the Postgres Database
	conn := db.ConnectDatabase(cfg.PostgresUrl)
	defer db.CloseConnection()

	// Initialize the User Profile Repository
	userProfile := repository.NewUserProfileRepository(ctx, conn)

	// Initialize the Mongo Database Client
	client := db.ConnectMongo(ctx, cfg.MongoUrl)
	defer db.CloseMongoClient(ctx)

	// Initialize the User Preferences Repository
	userPreferences := repository.NewUserPreferencesRepository(ctx, client)

	// Initialize Redis for session handling
	redis := cache.ConnectRedis(cfg.RedisUrl)
	defer redis.CloseRedisClient()

	// Initialize the Beer API Client
	beerApiClient := beerapi.NewClient(cfg.BeersApiUrl)

	// Initialize controller
	// NOTE: In production grade applications, the Controller should accept Service and Service would contain 1 or more
	// repositories, depending on the business logic implementation.
	ctrl := controller.New(beerApiClient, userProfile, userPreferences, redis)

	// The HTTP Server
	server.NewHttpServer(cfg.Port, router.Router(ctrl, cfg.StaticFiles)).
		WithContext(ctx).
		Run()
}
