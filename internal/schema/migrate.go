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
	quantity     INT,
	date_created TIMESTAMP,
	date_updated TIMESTAMP,

	PRIMARY KEY (card_id)
);`,
	},
// 	{
// 		Version:     2,
// 		Description: "Add sales",
// 		Script: `
// CREATE TABLE sales (
// 	sale_id      UUID,
// 	product_id   UUID,
// 	quantity     INT,
// 	paid         INT,
// 	date_created TIMESTAMP,

// 	PRIMARY KEY (sale_id),
// 	FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
// );`,
// 	},
	{
		Version:     2,
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
		Version:     3,
		Description: "Add user column to cards",
		Script: `
ALTER TABLE cards
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
