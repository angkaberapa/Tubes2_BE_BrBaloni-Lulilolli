package core

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// AppConfig holds all the application components
type AppContext struct {
	Config *AppConfig
	DBPool *pgxpool.Pool
}
