package schema

import (
	"github.com/GuiaBolso/darwin"
	"github.com/jmoiron/sqlx"
)

// migrations contains the queries needed to construct the database schema.
// Entries should never be removed from this slice once they have been ran in
// production.
//
// Including the queries directly in this file has the same pros/cons mentioned
// in seeds.go

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add cards",
		Script: `
CREATE TABLE cards (
	card_id   UUID,
	name         TEXT,
	content         TEXT,
	date_created TIMESTAMP,
	date_updated TIMESTAMP,

	PRIMARY KEY (card_id)
);`,
	},
	{
		Version:     2,
		Description: "Add menus",
		Script: `
CREATE TABLE menus (
	menu_id      UUID,
	name         TEXT,
	card_id   UUID,
	date_created TIMESTAMP,
	date_updated TIMESTAMP,

	PRIMARY KEY (menu_id),
	FOREIGN KEY (card_id) REFERENCES cards(card_id) ON DELETE CASCADE
);`,
	},
	{
		Version:     3,
		Description: "Add users",
		Script: `
CREATE TABLE users (
	user_id       UUID,
	name          TEXT,
	email         TEXT UNIQUE,
	roles         TEXT[],
	password_hash TEXT,

	date_created TIMESTAMP,
	date_updated TIMESTAMP,

	PRIMARY KEY (user_id)
);`,
	},
	{
		Version:     4,
		Description: "Add user column to cards",
		Script: `
ALTER TABLE cards
	ADD COLUMN user_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'
`,
	},
	{
		Version:     5,
		Description: "Add user column to menus",
		Script: `
ALTER TABLE menus
	ADD COLUMN user_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'
`,
	},
}

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {

	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}
