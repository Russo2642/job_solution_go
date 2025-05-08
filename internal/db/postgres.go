package db

import (
	"context"
	"database/sql"
	"fmt"

	"job_solition/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	db *sqlx.DB
}

func NewPostgreSQL(cfg config.PostgreSQLConfig) (*PostgreSQL, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к PostgreSQL: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка при проверке подключения к PostgreSQL: %w", err)
	}

	return &PostgreSQL{
		db: db,
	}, nil
}

func (p *PostgreSQL) Close() error {
	return p.db.Close()
}

func (p *PostgreSQL) Exec(query string, args ...interface{}) (sql.Result, error) {
	return p.db.Exec(query, args...)
}

func (p *PostgreSQL) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return p.db.ExecContext(ctx, query, args...)
}

func (p *PostgreSQL) Query(query string, args ...interface{}) (*sqlx.Rows, error) {
	return p.db.Queryx(query, args...)
}

func (p *PostgreSQL) QueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return p.db.QueryxContext(ctx, query, args...)
}

func (p *PostgreSQL) QueryRow(query string, args ...interface{}) *sqlx.Row {
	return p.db.QueryRowx(query, args...)
}

func (p *PostgreSQL) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return p.db.QueryRowxContext(ctx, query, args...)
}

func (p *PostgreSQL) Get(dest interface{}, query string, args ...interface{}) error {
	return p.db.Get(dest, query, args...)
}

func (p *PostgreSQL) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return p.db.GetContext(ctx, dest, query, args...)
}

func (p *PostgreSQL) Select(dest interface{}, query string, args ...interface{}) error {
	return p.db.Select(dest, query, args...)
}

func (p *PostgreSQL) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return p.db.SelectContext(ctx, dest, query, args...)
}

func (p *PostgreSQL) Begin() (*sqlx.Tx, error) {
	return p.db.Beginx()
}

func (p *PostgreSQL) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return p.db.BeginTxx(ctx, opts)
}

func (p *PostgreSQL) GetDB() *sqlx.DB {
	return p.db
}
