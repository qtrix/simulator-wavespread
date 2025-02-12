package db

import (
	"context"
)

// Get returns value for key
func (db *DB) Get(ctx context.Context, key string, out interface{}) error {
	sel := `
		select "value"
			from "store" s
			where "key" = $1;
		`
	row := db.QueryRow(ctx, sel, key)
	err := row.Scan(out)
	return err
}

// Set sets value for key, creates if not found
func (db *DB) Set(ctx context.Context, key string, value interface{}) error {
	upsert := `
	insert into "store" ("key", "value")
		values ($1, $2)
		on conflict ("key") do update set "value" = $2;
	`
	err := db.Exec(ctx, upsert, key, value)
	return err
}
