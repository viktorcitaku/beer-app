package config

import (
	"flag"
	"os"
)

var (
	port        string
	postgresUrl string
	mongoUrl    string
	redisUrl    string
	beersApiUrl string
	staticFiles string
)

func init() {
	flag.StringVar(&port, "port", "3000", "Application exposed port.")
	flag.StringVar(&postgresUrl, "postgres-url", "postgres://test:test@localhost:5432/beer", "The PostgresSQL URL.")
	flag.StringVar(&mongoUrl, "mongo-url", "mongodb://test:test@localhost:27017", "The Mongo DB URL")
	flag.StringVar(&redisUrl, "redis-url", "redis://localhost:6379", "The Redis URL")
	flag.StringVar(&beersApiUrl, "beers-api-url", "https://api.punkapi.com/v2/beers", "The Punk API URL")
	flag.Parse()

	// Fallback to environment variable
	tryEnvIfFlagNotSet(&port, "port", "PORT")
	tryEnvIfFlagNotSet(&postgresUrl, "postgres-url", "POSTGRES_URL")
	tryEnvIfFlagNotSet(&mongoUrl, "mongo-url", "MONGO_URL")
	tryEnvIfFlagNotSet(&redisUrl, "redis-url", "REDIS_URL")
	tryEnvIfFlagNotSet(&beersApiUrl, "beers-api-url", "BEERS_API_URL")
	tryEnvIfFlagNotSet(&staticFiles, "static-files", "BEER_STATIC_FILES")
}

func tryEnvIfFlagNotSet(cfg *string, flagName string, envName string) {
	if ok := isFlagSet(flagName); !ok {
		if env, ok := os.LookupEnv(envName); ok {
			*cfg = env
		}
	}
}

type Cfg struct {
	Port        string
	PostgresUrl string
	MongoUrl    string
	RedisUrl    string
	BeersApiUrl string
	StaticFiles string
}

func GetConfig() *Cfg {
	return &Cfg{
		Port:        port,
		PostgresUrl: postgresUrl,
		MongoUrl:    mongoUrl,
		RedisUrl:    redisUrl,
		BeersApiUrl: beersApiUrl,
		StaticFiles: staticFiles,
	}
}

func isFlagSet(flagName string) (found bool) {
	flag.Visit(func(f *flag.Flag) {
		if f.Name == flagName {
			found = true
		}
	})
	return
}
