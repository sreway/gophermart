package listener

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sreway/gophermart/pkg/postgres"
)

type PGListener struct {
	conn    *pgxpool.Conn
	channel string
}

func NewPGListener(ctx context.Context, pg *postgres.Postgres, channel string) (*PGListener, error) {
	conn, err := pg.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	_, err = conn.Exec(ctx, "LISTEN "+channel)
	if err != nil {
		return nil, err
	}
	return &PGListener{
		conn:    conn,
		channel: channel,
	}, nil
}

func (pgl *PGListener) Listen(ctx context.Context, data chan<- string) {
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

func (pgl *PGListener) Release() {
	pgl.conn.Release()
}
