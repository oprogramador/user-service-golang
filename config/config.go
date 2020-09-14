package config

import (
	"time"
)

func GetMongoURL() string {
	return "mongodb://localhost:27017"
}

func GetTimeout() time.Duration {
	return 120 * time.Second
}

func GetPort() string {
	return "10000"
}
