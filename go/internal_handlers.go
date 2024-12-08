package main

import (
	"database/sql"
	"errors"
	"net/http"
)

// このAPIをインスタンス内から一定間隔で叩かせることで、椅子とライドをマッチングさせる
func internalGetMatching(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// MEMO: 一旦最も待たせているリクエストに適当な空いている椅子マッチさせる実装とする。おそらくもっといい方法があるはず…
	var rides []Ride
	if err := db.SelectContext(ctx, &rides, `SELECT * FROM rides WHERE chair_id IS NULL ORDER BY created_at`); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	var chairs []Chair

	if err := db.SelectContext(ctx, &chairs, "SELECT * FROM chairs WHERE is_active = TRUE"); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if len(chairs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	matchedRides := make(map[string]bool)

	for i := range chairs {
		var empty bool
		if err := db.GetContext(ctx, &empty, "SELECT COUNT(*) = 0 FROM (SELECT COUNT(chair_sent_at) = 6 AS completed FROM ride_statuses WHERE ride_id IN (SELECT id FROM rides WHERE chair_id = ?) GROUP BY ride_id) is_completed WHERE completed = FALSE", chairs[i].ID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		if empty {
			break
		}

		var chairLocation ChairLocation

		if err := db.SelectContext(ctx, &chairLocation, "SELECT * FROM chair_locations WHERE chair_id = ? ORDER BY created_at DESC LIMIT 1", chairs[i].ID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		mind := 1 << 31
		minj := 0
		for j := range rides {
			if _, ok := matchedRides[rides[j].ID]; ok {
				continue
			}
			dist := abs(chairLocation.Latitude-rides[j].DestinationLatitude) + abs(chairLocation.Longitude-rides[j].DestinationLongitude)
			if dist < mind {
				mind = dist
				minj = j
			}
		}

		if _, err := db.ExecContext(ctx, "UPDATE rides SET chair_id = ? WHERE id = ?", chairs[i].ID, rides[minj].ID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		matchedRides[rides[minj].ID] = true
	}

	matched := &Chair{}
	empty := false
	for i := 0; i < 10; i++ {
		if err := db.GetContext(ctx, matched, "SELECT * FROM chairs INNER JOIN (SELECT id FROM chairs WHERE is_active = TRUE ORDER BY RAND() LIMIT 1) AS tmp ON chairs.id = tmp.id LIMIT 1"); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			writeError(w, http.StatusInternalServerError, err)
		}

		if err := db.GetContext(ctx, &empty, "SELECT COUNT(*) = 0 FROM (SELECT COUNT(chair_sent_at) = 6 AS completed FROM ride_statuses WHERE ride_id IN (SELECT id FROM rides WHERE chair_id = ?) GROUP BY ride_id) is_completed WHERE completed = FALSE", matched.ID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		if empty {
			break
		}
	}
	if !empty {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if _, err := db.ExecContext(ctx, "UPDATE rides SET chair_id = ? WHERE id = ?", matched.ID, ride.ID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
