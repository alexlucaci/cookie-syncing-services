package schema

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}

// migrations contains the queries needed to construct the database schema.
// Entries should never be removed from this slice once they have been ran in
// production.
var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add users",
		Script: `
		CREATE TABLE cookies (
			cookie_id UUID NOT NULL,
			partner_cookie_id UUID,
			date_created TIMESTAMP,
			date_updated TIMESTAMP,
	
			PRIMARY KEY (cookie_id)
		);`,
	},
}
