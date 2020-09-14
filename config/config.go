package config

import (
	"os"
	"time"
)

func GetMongoURL() string {
	env := os.Getenv("MONGO_URL")
	if env != "" {
		return env
	}
	return "mongodb://localhost:27017"
}

func GetTimeout() (time.Duration, error) {
	env := os.Getenv("TIMEOUT")
	if env != "" {
		timeout, err := time.ParseDuration(env)
		if err != nil {
			return 0, err
		}
		return timeout, nil
	}
	return 120 * time.Second, nil
}

func GetPort() string {
	env := os.Getenv("PORT")
	if env != "" {
		return env
	}
	return "10000"
}
