package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

const (
	ScmshPortKey          = "SCMSH_PORT"
	ScmshRedisEnabledKey  = "SCMSH_REDIS_ENABLED"
	ScmshRedisAddressKey  = "SCMSH_REDIS_ADDRESS"
	ScmshRedisPasswordKey = "SCMSH_REDIS_PASSWORD"
	ScmshRedisDatabaseKey = "SCMSH_REDIS_DATABASE"
)

type Config struct {
	Port          int
	RedisEnabled  bool
	RedisAddress  string
	RedisPassword string
	RedisDatabase int
}

func Get() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using system env vars")
	}

	port := 8080
	if portSTR := os.Getenv(ScmshPortKey); portSTR != "" {
		port, err = strconv.Atoi(portSTR)
		if err != nil {
			return Config{}, errors.Wrapf(err, "invalid port value %q", portSTR)
		}
	}

	redisEnabled := false
	if enabledStr := os.Getenv(ScmshRedisEnabledKey); enabledStr != "" {
		redisEnabled, err = strconv.ParseBool(enabledStr)
		if err != nil {
			return Config{}, errors.Wrapf(err, "invalid redis enabled value %q", enabledStr)
		}
	}

	redisAddress := "localhost:6379"
	if address := os.Getenv(ScmshRedisAddressKey); address != "" {
		redisAddress = address
	}

	redisPassword := ""
	if password := os.Getenv(ScmshRedisPasswordKey); password != "" {
		redisPassword = password
	}

	redisDatabase := 0
	if database := os.Getenv(ScmshRedisDatabaseKey); database != "" {
		redisDatabase, err = strconv.Atoi(database)
		if err != nil {
			return Config{}, errors.Wrapf(err, "invalid redis database value %q", database)
		}
	}

	return Config{
		Port:          port,
		RedisEnabled:  redisEnabled,
		RedisAddress:  redisAddress,
		RedisPassword: redisPassword,
		RedisDatabase: redisDatabase,
	}, nil
}
