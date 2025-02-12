package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTablePriceData, nil)
}

func upCreateTablePriceData(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create table price_data
		(
			timestamp timestamp with time zone not null
				constraint price_data_pk primary key,
			price     numeric(10,2)
		);
	`)
	return err
}
