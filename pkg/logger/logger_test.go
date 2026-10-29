package logger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitDevelopment(t *testing.T) {
	log := Init("development")
	assert.NotNil(t, log)
	assert.Equal(t, log, Log)
	assert.Equal(t, log, zap.L())
	assert.False(t, IsProduction())

	log.Info("test development log message", zap.String("key", "val"))
	Sync()
}

func TestInitProduction(t *testing.T) {
	log := Init("production")
	assert.NotNil(t, log)
	assert.Equal(t, log, Log)
	assert.Equal(t, log, zap.L())
	assert.True(t, IsProduction())

	log.Info("test production log message", zap.String("key", "val"))
	Sync()
}

func TestInitFallbackEnv(t *testing.T) {
	os.Setenv("APP_ENV", "production")
	defer os.Unsetenv("APP_ENV")

	log := Init()
	assert.NotNil(t, log)
	assert.True(t, IsProduction())
}
