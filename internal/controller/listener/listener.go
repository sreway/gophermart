package listener

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sreway/gophermart/pkg/postgres"
)

type PgListener struct {
	conn    *pgxpool.Conn
	channel string
}

func NewPgListner(ctx context.Context, pg *postgres.Postgres, channel string) (*PgListener, error) {
	conn, err := pg.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	_, err = conn.Exec(ctx, "LISTEN "+channel)
	if err != nil {
		return nil, err
	}
	return &PgListener{
		conn:    conn,
		channel: channel,
	}, nil
}

func (pgl *PgListener) Listen(ctx context.Context, data chan<- string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			notify, errNotify := pgl.conn.Conn().WaitForNotification(ctx)
			if errNotify != nil {
				time.Sleep(10 * time.Millisecond)
			}

			if ctx.Err() == nil {
				data <- notify.Payload
			}
		}
	}
}

func (pgl *PgListener) Release() {
	pgl.conn.Release()
}
