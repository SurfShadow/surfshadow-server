package db

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib" // Import for side effects
	"github.com/jmoiron/sqlx"

	"github.com/SurfShadow/surfshadow-server/internal/infrastructure/config"
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

	db, err := sqlx.Connect(c.PgDriver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
