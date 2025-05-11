package postgres

import (
	"database/sql"
	"errors"
	"ndx/pkg/logger"

	"fmt"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"net/url"
)

type Config struct {
	Host         string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" env-default:"localhost"`
	Port         uint16 `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" env-default:"5432"`
	Username     string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" env-default:"root"`
	Password     string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD" env-default:"password"`
	Database     string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" env-default:"postgres"`
	MinConns     int32  `yaml:"POSTGRES_MIN_CONNS" env:"POSTGRES_MIN_CONNS" env-default:"5"`
	MaxConns     int32  `yaml:"POSTGRES_MAX_CONNS" env:"POSTGRES_MAX_CONNS" env-default:"10"`
	SearchSchema string `yaml:"POSTGRES_MAIN_SCHEMA" env:"POSTGRES_MAIN_SCHEMA" env-default:"base_schema"`
}

func New(config Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.L().Logf(0, "can't open db | err: %v", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		logger.L().Logf(0, "can't ping db | err: %v", err)
		return nil, err
	}

	user := url.QueryEscape(config.Username)
	pass := url.QueryEscape(config.Password)
	schema := url.QueryEscape(config.SearchSchema)

	connStr2 := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s",
		user,
		pass,
		config.Host,
		config.Port,
		config.Database,
		schema,
	)

	m, err := migrate.New("file://db/migrations", connStr2)
	if err != nil {
		logger.L().Logf(0, "can't migrat6e db | err: %v", err)
		return nil, err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.L().Logf(0, "can't migrate 2.0 db | err: %v", err)
		return db, nil
	}

	return db, nil
}
