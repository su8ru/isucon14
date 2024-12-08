package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/motoki317/sc"
	"time"
)

type DBTX interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

var chairCache, _ = sc.New(func(ctx context.Context, chairID string) (Chair, error) {
	tx, ok := ctx.Value("tx").(DBTX)
	if !ok {
		tx = db
	}

	var chair Chair
	if err := tx.GetContext(ctx, &chair, "SELECT * FROM chairs WHERE id = ?", chairID); err != nil {
		return Chair{}, err
	}

	return chair, nil
}, 1*time.Minute, 5*time.Minute)

var chairFromAccessTokenCache, _ = sc.New(func(ctx context.Context, accessToken string) (Chair, error) {
	tx, ok := ctx.Value("tx").(DBTX)
	if !ok {
		tx = db
	}

	var chair Chair
	if err := tx.GetContext(ctx, &chair, "SELECT * FROM chairs WHERE access_token = ?", accessToken); err != nil {
		return Chair{}, err
	}

	return chair, nil
}, 1*time.Minute, 5*time.Minute)

var latestRideStatusCache, _ = sc.New(func(ctx context.Context, rideID string) (string, error) {
	tx, ok := ctx.Value("tx").(DBTX)
	if !ok {
		tx = db
	}

	status := ""
	if err := tx.GetContext(ctx, &status, `SELECT status FROM ride_statuses WHERE ride_id = ? ORDER BY created_at DESC LIMIT 1`, rideID); err != nil {
		return "", err
	}
	return status, nil
}, 1*time.Minute, 5*time.Minute)

var yetChairSentRideStatusCache, _ = sc.New(func(ctx context.Context, rideID string) (RideStatus, error) {
	tx, ok := ctx.Value("tx").(DBTX)
	if !ok {
		tx = db
	}

	var rideStatus RideStatus
	if err := tx.GetContext(ctx, &rideStatus, `SELECT * FROM ride_statuses WHERE ride_id = ? AND chair_sent_at IS NULL ORDER BY created_at ASC LIMIT 1`, rideID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return RideStatus{}, nil
		}
		return RideStatus{}, err
	}
	return rideStatus, nil
}, 1*time.Minute, 5*time.Minute)
