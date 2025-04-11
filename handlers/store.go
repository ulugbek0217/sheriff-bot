package handlers

import (
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/ulugbek0217/sheriff-bot/db/sqlc"
)

type App struct {
	Store  db.Store
	Pool   *pgxpool.Pool
	Admins []int64
	WG     *sync.WaitGroup
}
