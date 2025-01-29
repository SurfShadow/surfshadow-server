package db

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"

	"github.com/SurfShadow/surfshadow-server/internal/infrastructure/config"
	"github.com/SurfShadow/surfshadow-server/pkg/logger"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

func NewPsqlDB(c *config.DBConfig) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
		c.Host,
		c.Port,
		c.User,
		c.DataBaseName,
		c.Password,
	)

	logger.Instance.Debugf("DSN: %s", dataSourceName)

	logger.Instance.Info("Connecting to database")

	db, err := sqlx.Connect(c.PgDriver, dataSourceName)
	if err != nil {
		logger.Instance.Errorf("Failed to connect to database: %v", err)
		return nil, err
	}

	logger.Instance.Info("Database connected successfully")

	logger.Instance.Debugf("Max open connections: %d", maxOpenConns)
	db.SetMaxOpenConns(maxOpenConns)

	logger.Instance.Debugf("Connection max lifetime: %d", connMaxLifetime)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)

	logger.Instance.Debugf("Max idle connections: %d", maxIdleConns)
	db.SetMaxIdleConns(maxIdleConns)

	logger.Instance.Debugf("Connection max idle time: %d", connMaxIdleTime)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	logger.Instance.Info("Pinging database to verify connection")

	if err = db.Ping(); err != nil {
		logger.Instance.Errorf("Failed to ping database: %v", err)
		return nil, err
	}

	logger.Instance.Info("Database ping successful")

	return db, nil
}
