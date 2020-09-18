package cookie

import "time"

// We will use pointer to primitive types because it's value in db
// can be null and we have to unmarshall it
type Cookie struct {
	ID              string    `db:"cookie_id"`
	PartnerCookieID *string   `db:"partner_cookie_id"`
	DateCreated     time.Time `db:"date_created"`
	DateUpdated     time.Time `db:"date_updated"`
}
