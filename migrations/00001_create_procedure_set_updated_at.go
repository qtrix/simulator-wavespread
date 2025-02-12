package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateSetUpdatedAt, nil)
}

func upCreateSetUpdatedAt(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create or replace function trigger_set_updated_at() returns TRIGGER as
		$$
		begin
			NEW.updated_at = now();
			return NEW;
		end;
		$$ language plpgsql;
	`)
	return err
}
