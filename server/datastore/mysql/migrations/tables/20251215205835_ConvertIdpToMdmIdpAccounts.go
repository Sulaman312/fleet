package tables

import (
	"database/sql"

	"github.com/pkg/errors"
)

func init() {
	MigrationClient.AddMigration(Up_20251215205835, Down_20251215205835)
}

func Up_20251215205835(tx *sql.Tx) error {
	// https://github.com/fleetdm/fleet/issues/37168
	_, err := tx.Exec(`UPDATE host_emails SET source = 'mdm_idp_accounts' WHERE source = 'idp'`)
	if err != nil {
		return errors.Wrap(err, "update idp source to mdm_idp_accounts")
	}
	return nil
}

func Down_20251215205835(tx *sql.Tx) error {
	return nil
}
