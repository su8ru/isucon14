package main

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

func updateRideStatus(tx *sqlx.Tx, ctx context.Context, id string, rideId string, status string) error {
	_, err := tx.ExecContext(ctx, "INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)", ulid.Make().String(), rideId, status)
	return err
}
