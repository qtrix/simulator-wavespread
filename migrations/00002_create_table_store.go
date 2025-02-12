package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableStore, nil)
}

func upCreateTableStore(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create table "store"
		(
			"key"        text not null
				constraint "store_pk" primary key,
			"value"      text,
			"created_at" timestamp default now(),
			"updated_at" timestamp default now()
		);
		
		create unique index "store_key_uindex" on "store" ("key");
		
		create trigger "set_timestamp"
			before update
			on "store"
			for each row
		execute procedure trigger_set_updated_at();
	`)
	return err
}
