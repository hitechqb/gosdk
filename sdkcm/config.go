package sdkcm

import (
	"os"
	"strconv"
	"time"
)

var (
	cf SDKConfig
)

func SetConfig(c SDKConfig) {
	cf = c
}

func GetConfig() SDKConfig {
	return cf
}

type SDKConfig interface {
	// App
	AppPort() string
	AppURLScheme() string

	// Database connection
	DbHost() string
	DbPort() string
	DbUser() string
	DbPassword() string
	DbName() string

	NatURL() string
	RedisURL() string

	RequestTimeout() time.Duration
	CacheLifetime() time.Duration
	ItemsPerPage() int
}

type config struct {
}

func NewSDKConfig() *config {
	return &config{}
}

func (s config) AppPort() string {
	return os.Getenv("APP_PORT")
}

func (s config) AppURLScheme() string {
	return os.Getenv("APP_URL_SCHEME")
}

func (s *config) NatURL() string {
	return os.Getenv("NAT_URL")
}

func (s config) RequestTimeout() time.Duration {
	second, _ := strconv.Atoi(os.Getenv("REQUEST_TIMEOUT"))
	return time.Duration(second) * time.Second
}

func (s config) DbHost() string {
	return os.Getenv("DB_HOST")
}

func (s config) DbPort() string {
	return os.Getenv("DB_PORT")
}

func (s config) DbUser() string {
	return os.Getenv("DB_USER")
}

func (s config) DbPassword() string {
	return os.Getenv("DB_PASSWORD")
}

func (s config) DbName() string {
	return os.Getenv("DB_NAME")
}

func (s config) RedisURL() string {
	return os.Getenv("REDIS_URL")
}

func (s config) CacheLifetime() time.Duration {
	minute, _ := strconv.Atoi(os.Getenv("CACHE_LIFETIME"))
	return time.Duration(minute) * time.Second
}

func (s *config) ItemsPerPage() int {
	x, _ := strconv.Atoi(os.Getenv("ITEMS_PER_PAGE"))
	if x <= 0 {
		return 5
	}

	return x
}
