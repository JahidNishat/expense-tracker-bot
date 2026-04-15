package configs

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/modules/cache"
	"github.com/masudur-rahman/expense-tracker-bot/modules/google"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"github.com/masudur-rahman/styx"
	isql "github.com/masudur-rahman/styx/sql"
	"github.com/masudur-rahman/styx/sql/postgres"
	"github.com/masudur-rahman/styx/sql/sqlite"
	"github.com/masudur-rahman/styx/sql/sqlite/lib"

	_ "github.com/lib/pq"
)

// sqlDB holds a reference to the database engine for utility functions.
var (
	sqlDB isql.Engine
	dbMu  sync.Mutex
)

// GetUnitOfWork returns a UnitOfWork wrapping the active database engine.
func GetUnitOfWork() styx.UnitOfWork {
	return styx.UnitOfWork{SQL: &safeEngine{engine: sqlDB, mu: &dbMu}}
}

func InitiateCache() {
	cache.Init(TrackerConfig.Cache)
}

func InitiateDatabaseConnection(ctx context.Context) error {
	cfg := TrackerConfig.Database
	switch cfg.Type {
	case DatabasePostgres:
		db, err := getPostgresDatabase(ctx)
		if err != nil {
			return err
		}
		return initializeSQLServices(styx.UnitOfWork{SQL: db})
	case DatabaseSQLite, "":
		if cfg.SQLite.SyncToDrive {
			if !cfg.SQLite.DisableSyncFromDrive {
				if err := google.SyncDatabaseFromDrive(); err != nil {
					return err
				}
				logr.DefaultLogger.Infof("SQLite database synced from google drive")
			}
			go google.SyncDatabaseToDrivePeriodically(TrackerConfig.Database.SQLite.SyncInterval)
		}

		db, err := getSQLiteDatabase(ctx)
		if err != nil {
			return err
		}
		return initializeSQLServices(styx.UnitOfWork{SQL: db})
	default:
		return fmt.Errorf("unknown database type")
	}
}

func getSQLiteDatabase(ctx context.Context) (isql.Engine, error) {
	conn, err := lib.GetSQLiteConnection(google.DatabasePath())
	if err != nil {
		return nil, err
	}

	return sqlite.NewSQLite(ctx, conn), nil
}

func initializeSQLServices(uow styx.UnitOfWork) error {
	dbMu.Lock()
	sqlDB = uow.SQL
	dbMu.Unlock()

	// Use safeEngine for services to ensure thread safety
	safe := &safeEngine{engine: uow.SQL, mu: &dbMu}
	uow.SQL = safe

	if err := syncTables(uow.SQL); err != nil {
		return err
	}
	if err := fixNullZeroValues(uow.SQL); err != nil {
		return err
	}
	all.InitiateSQLServices(uow, logr.DefaultLogger)

	return all.GetServices().Txn.UpdateTxnCategories()
}

// fixNullZeroValues patches existing rows where styx v1.2.3 inserted NULL for zero-value fields.
func fixNullZeroValues(db isql.Engine) error {
	stmts := []string{
		`UPDATE "transaction" SET deleted_at = 0 WHERE deleted_at IS NULL`,
		`UPDATE "wallet" SET version = 0 WHERE version IS NULL`,
		`UPDATE "contacts" SET net_balance = 0 WHERE net_balance IS NULL`,
		`UPDATE "contacts" SET last_txn_timestamp = 0 WHERE last_txn_timestamp IS NULL`,
		`UPDATE "budget" SET alert_at = 80 WHERE alert_at IS NULL`,
		`UPDATE "refresh_token" SET revoked = 0 WHERE revoked IS NULL`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("fix null values: %w", err)
		}
	}
	return nil
}

// pgPool holds the *sql.DB connection pool so we can re-acquire connections.
var (
	pgPool     *sql.DB
	activeConn *sql.Conn
)

func getPostgresDatabase(ctx context.Context) (isql.Engine, error) {
	parsePostgresConfig()

	pool, err := sql.Open("postgres", TrackerConfig.Database.Postgres.String())
	if err != nil {
		return nil, fmt.Errorf("open postgres pool: %w", err)
	}

	pool.SetMaxOpenConns(25)
	pool.SetMaxIdleConns(5)
	pool.SetConnMaxLifetime(5 * time.Minute)

	if err := pool.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	conn, err := pool.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquire initial conn: %w", err)
	}

	pgPool = pool
	activeConn = conn
	go pingPostgresDatabasePeriodically(context.Background(), logr.DefaultLogger)

	return postgres.NewPostgres(ctx, conn).ShowSQL(true), nil
}

func pingPostgresDatabasePeriodically(ctx context.Context, logger logr.Logger) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if pgPool == nil {
			continue
		}

		if err := PingDatabase(); err != nil {
			logger.Warnw("Postgres connection lost, attempting to re-acquire...", "error", err.Error())

			conn, connErr := pgPool.Conn(ctx)
			if connErr != nil {
				logger.Errorw("Critical: Failed to acquire fresh connection from pool", "error", connErr.Error())
				continue
			}

			if activeConn != nil {
				_ = activeConn.Close()
			}

			activeConn = conn
			newEngine := postgres.NewPostgres(ctx, conn).ShowSQL(true)

			dbMu.Lock()
			sqlDB = newEngine
			dbMu.Unlock()

			all.InitiateSQLServices(styx.UnitOfWork{SQL: &safeEngine{engine: newEngine, mu: &dbMu}}, logger)
			logger.Infow("Successfully re-established Postgres connection")
		}
	}
}

func parsePostgresConfig() {
	user, ok := os.LookupEnv("POSTGRES_USER")
	if ok {
		TrackerConfig.Database.Postgres.User = user
	}
	pass, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if ok {
		TrackerConfig.Database.Postgres.Password = pass
	}
	name, ok := os.LookupEnv("POSTGRES_DB")
	if ok {
		TrackerConfig.Database.Postgres.Name = name
	}
	host, ok := os.LookupEnv("POSTGRES_HOST")
	if ok {
		TrackerConfig.Database.Postgres.Host = host
	}
	port, ok := os.LookupEnv("POSTGRES_PORT")
	if ok {
		TrackerConfig.Database.Postgres.Port = port
	}
	ssl, ok := os.LookupEnv("POSTGRES_SSL_MODE")
	if ok {
		TrackerConfig.Database.Postgres.SSLMode = ssl
	}
}

func syncTables(db isql.Engine) error {
	return db.Sync(
		models.Profile{},
		models.Contacts{},
		models.Wallet{},
		models.Transaction{},
		models.TxnCategory{},
		models.TxnSubcategory{},
		models.Event{},
		models.AICache{},
		models.Budget{},
		models.RefreshToken{},
	)
}

// LoadAICacheIntoMemory loads all persisted AI cache rows into the in-memory cache.
func LoadAICacheIntoMemory() {
	if err := PingDatabase(); err != nil {
		return
	}

	var rows []models.AICache
	if err := GetUnitOfWork().SQL.Table(models.AICache{}.TableName()).FindMany(&rows); err != nil {
		logr.DefaultLogger.Errorw("Failed to load AI cache", "error", err.Error())
		return
	}
	for _, row := range rows {
		_ = cache.SetCache(row.InputText, row.SubcategoryID, -1)
	}
	logr.DefaultLogger.Infow("AI cache loaded from DB", "count", len(rows))
}

// PingDatabase checks if the database connection is healthy.
func PingDatabase() error {
	dbMu.Lock()
	defer dbMu.Unlock()
	if sqlDB == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := sqlDB.Exec("SELECT 1")
	return err
}

// InsertAICache persists a single AI cache entry to the database.
func InsertAICache(entry models.AICache) error {
	if sqlDB == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := GetUnitOfWork().SQL.Table(models.AICache{}.TableName()).InsertOne(entry)
	return err
}

// safeEngine is a thread-safe wrapper around isql.Engine.
type safeEngine struct {
	engine isql.Engine
	mu     *sync.Mutex
}

func (s *safeEngine) BeginTx() (isql.Engine, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	engine, err := s.engine.BeginTx()
	if err != nil {
		return nil, err
	}
	return &safeEngine{engine: engine, mu: s.mu}, nil
}

func (s *safeEngine) Commit() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.engine.Commit()
}

func (s *safeEngine) Rollback() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.engine.Rollback()
}

func (s *safeEngine) Table(name string) isql.Engine {
	return &safeEngine{engine: s.engine.Table(name), mu: s.mu}
}

func (s *safeEngine) ID(id any) isql.Engine {
	return &safeEngine{engine: s.engine.ID(id), mu: s.mu}
}

func (s *safeEngine) In(col string, values ...any) isql.Engine {
	return &safeEngine{engine: s.engine.In(col, values...), mu: s.mu}
}

func (s *safeEngine) Where(cond string, args ...any) isql.Engine {
	return &safeEngine{engine: s.engine.Where(cond, args...), mu: s.mu}
}

func (s *safeEngine) Columns(cols ...string) isql.Engine {
	return &safeEngine{engine: s.engine.Columns(cols...), mu: s.mu}
}

func (s *safeEngine) AllCols() isql.Engine {
	return &safeEngine{engine: s.engine.AllCols(), mu: s.mu}
}

func (s *safeEngine) MustCols(cols ...string) isql.Engine {
	return &safeEngine{engine: s.engine.MustCols(cols...), mu: s.mu}
}

func (s *safeEngine) ShowSQL(showSQL bool) isql.Engine {
	return &safeEngine{engine: s.engine.ShowSQL(showSQL), mu: s.mu}
}

func (s *safeEngine) FindOne(document any, filter ...any) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.engine.FindOne(document, filter...)
}

func (s *safeEngine) FindMany(documents any, filter ...any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.engine.FindMany(documents, filter...)
}

func (s *safeEngine) InsertOne(document any) (any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.engine.InsertOne(document)
}

func (s *safeEngine) InsertMany(documents []any) ([]any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.engine.InsertMany(documents)
}

func (s *safeEngine) UpdateOne(document any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.engine.UpdateOne(document)
}

func (s *safeEngine) DeleteOne(filter ...any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.engine.DeleteOne(filter...)
}

func (s *safeEngine) Query(query string, args ...any) (*sql.Rows, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.engine.Query(query, args...)
}

func (s *safeEngine) Exec(query string, args ...any) (sql.Result, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.engine.Exec(query, args...)
}

func (s *safeEngine) Sync(tables ...any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.engine.Sync(tables...)
}

func (s *safeEngine) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.engine.Close()
}
