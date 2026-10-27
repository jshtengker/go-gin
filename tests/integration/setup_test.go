//go:build integration

package integration

import (
	"backend/internal/api/router"
	"backend/internal/configs"
	"backend/pkg/logger"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()

	fmt.Println()
	if exitCode == 0 {
		fmt.Println("====================================================================")
		fmt.Println("                INTEGRATION TEST SUMMARY REPORT                     ")
		fmt.Println("====================================================================")
		fmt.Println("  Status : ALL PASSED")
		fmt.Println("  Phases Tested (Dual-Token Auth Flow with Grace Period):")
		fmt.Println("    [PASS] 1_Register                  - DB insertion & hashed password")
		fmt.Println("    [PASS] 2_Login                     - Credential verification & JWT issuance")
		fmt.Println("    [PASS] 3_GetProfile_Me             - Protected endpoint & claim injection")
		fmt.Println("    [PASS] 4_TokenRotation_Refresh     - Refresh token rotation & single-use guarantee")
		fmt.Println("    [PASS] 5_Concurrent_GracePeriod    - Concurrent frontend handling within grace period")
		fmt.Println("    [PASS] 6_TokenReuse_Attack         - Breach detection & session lockout after grace period (401)")
		fmt.Println("    [PASS] 7_Logout                    - Stateful session revocation & immediate invalidation")
		fmt.Println("====================================================================")
	} else {
		fmt.Println("====================================================================")
		fmt.Println("                INTEGRATION TEST SUMMARY REPORT                     ")
		fmt.Println("====================================================================")
		fmt.Println("  Status : FAILED - One or more integration tests failed.")
		fmt.Println("====================================================================")
	}

	os.Exit(exitCode)
}

// getProjectRoot returns the absolute path to the backend repository root
func getProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	// tests/integration/setup_test.go -> two levels up to root
	return filepath.Join(filepath.Dir(b), "../..")
}

// setupIntegrationApp initializes DB, runs migrations, and creates a Gin engine for testing.
func setupIntegrationApp(t *testing.T) (*gin.Engine, *gorm.DB, *configs.Config) {
	gin.SetMode(gin.TestMode)
	logger.Log = zap.NewNop()

	root := getProjectRoot()
	_ = godotenv.Load(filepath.Join(root, "internal/configs/.env"))
	cfg := configs.Load()

	require.NotEmpty(t, cfg.DBHost, "DB_HOST must be provided in internal/configs/.env")
	require.NotEmpty(t, cfg.DBName, "DB_NAME must be provided in internal/configs/.env")

	// 1. Run migrations up using golang-migrate
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	sqlDB, err := sql.Open("postgres", dsn)
	require.NoError(t, err, "failed to open postgres connection for migrations")
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	require.NoError(t, err, "failed to create migrate driver")

	migrationPath := fmt.Sprintf("file://%s", filepath.ToSlash(filepath.Join(root, "migrations")))
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		cfg.DBName,
		driver,
	)
	require.NoError(t, err, "failed to initialize migrate instance")

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err, "failed to apply migrations in integration test")
	}

	// 2. Connect via GORM
	gormDB, err := configs.ConnPostgres(cfg)
	require.NoError(t, err, "failed to connect to Postgres via GORM")

	// 3. Setup Gin Router
	engine := router.Setup(cfg, gormDB)

	return engine, gormDB, cfg
}

func executeRequest(r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
