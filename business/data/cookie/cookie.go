package cookie

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

var (
	// ErrNotFound is used when a specific Cookie is requested but not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Create will add a new Cookie without the PartnerCookieID in the db
func Create(ctx context.Context, db *sqlx.DB, now time.Time) (Cookie, error) {
	c := Cookie{
		ID:          uuid.New().String(),
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `INSERT into cookies (cookie_id, date_created, date_updated) VALUES($1, $2, $3)`
	if _, err := db.ExecContext(ctx, q, c.ID, c.DateCreated, c.DateUpdated); err != nil {
		return Cookie{}, errors.Wrap(err, "inserting cookie into db")
	}

	return c, nil
}

// SetPartnerCookie will set PartnerCookieID for an existing Cookie
func SetPartnerCookie(ctx context.Context, db *sqlx.DB, id, partnerCookieId string, now time.Time) error {
	// First check if the ids are valid uuids
	for _, cID := range []string{id, partnerCookieId} {
		if _, err := uuid.Parse(cID); err != nil {
			return ErrInvalidID
		}
	}

	const q = `UPDATE cookies SET partner_cookie_id = $2, date_updated = $3 WHERE cookie_id = $1`

	if _, err := db.ExecContext(ctx, q, id, partnerCookieId, now); err != nil {
		return errors.Wrap(err, "setting partner cookie")
	}

	return nil
}

// One is retrieving a cookie by its id
func One(ctx context.Context, db *sqlx.DB, id string) (Cookie, error) {
	if _, err := uuid.Parse(id); err != nil {
		return Cookie{}, ErrInvalidID
	}

	const q = `SELECT * from cookies WHERE cookie_id = $1`

	var c Cookie
	if err := db.GetContext(ctx, &c, q, id); err != nil {
		if err == sql.ErrNoRows {
			return Cookie{}, ErrNotFound
		}

		return Cookie{}, errors.Wrap(err, "selecting cookie from db")
	}

	return c, nil
}
