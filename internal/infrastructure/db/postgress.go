package db

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/SurfShadow/surfshadow-server/internal/infrastructure/config"
	"time"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

func NewPsqlDB(c *config.Config) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.DB.Host,
		c.DB.Port,
		c.DB.User,
		c.DB.DataBaseName,
		c.DB.Password,
	)

	db, err := sqlx.Connect(c.DB.PgDriver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
