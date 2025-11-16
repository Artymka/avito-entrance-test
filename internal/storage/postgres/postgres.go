package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"

	"github.com/Artymka/avito-entrance-test/internal/config"
)

type Postgres struct {
	DB *sqlx.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	const op = "postgres.new"
	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf(
			"dbname=%s host=%s port=%s sslmode=%s user=%s password=%s",
			cfg.Postgres.Database,
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.SSLMode,
			cfg.Postgres.Username,
			cfg.Postgres.Password,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Postgres{DB: db}, nil
}

func (p *Postgres) Close() error {
	const op = "postgres.close"
	err := p.DB.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
