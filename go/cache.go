package main

import (
	"context"
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
