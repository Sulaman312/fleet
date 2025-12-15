package tables

import (
	"testing"

	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/stretchr/testify/require"
)

func TestUp_20251215205835(t *testing.T) {
	db := applyUpToPrev(t)

	hostId := execNoErrLastID(t, db, //nolint:gosec // dismiss G115
		`INSERT INTO hosts (osquery_host_id, node_key, uuid, platform) VALUES (?, ?, ?, ?);`,
		1, 1, "macOS_UUID", "darwin",
	)

	execNoErr(t, db,
		`INSERT INTO host_emails (email, host_id, source) VALUES (?, ?, ?)`,
		"mdm.user.idp@example.com", hostId, "idp",
	)

	execNoErr(t, db,
		`INSERT INTO host_emails (email, host_id, source) VALUES (?, ?, ?)`,
		"mdm.user.chrome@example.com", hostId, fleet.DeviceMappingGoogleChromeProfiles,
	)

	execNoErr(t, db,
		`INSERT INTO host_emails (email, host_id, source) VALUES (?, ?, ?)`,
		"mdm.user.mdm.idp@example.com", hostId, fleet.DeviceMappingMDMIdpAccounts,
	)

	execNoErr(t, db,
		`INSERT INTO host_emails (email, host_id, source) VALUES (?, ?, ?)`,
		"mdm.user.custom.installer@example.com", hostId, fleet.DeviceMappingCustomInstaller,
	)
	execNoErr(t, db,
		`INSERT INTO host_emails (email, host_id, source) VALUES (?, ?, ?)`,
		"mdm.user.custom.override@example.com", hostId, fleet.DeviceMappingCustomOverride,
	)
	execNoErr(t, db,
		`INSERT INTO host_emails (email, host_id, source) VALUES (?, ?, ?)`,
		"mdm.user.custom.replacement@example.com", hostId, fleet.DeviceMappingCustomReplacement,
	)

	applyNext(t, db)

	var hostEmails []struct {
		Email  string `db:"email"`
		Source string `db:"source"`
	}

	err := db.Select(&hostEmails, `SELECT email, source FROM host_emails WHERE host_id = ?;`, hostId)
	require.NoError(t, err)
	require.Len(t, hostEmails, 6)

	for _, he := range hostEmails {
		switch he.Email {
		case "mdm.user.idp@example.com":
			require.Equal(t, fleet.DeviceMappingMDMIdpAccounts, he.Source)
		case "mdm.user.chrome@example.com":
			require.Equal(t, fleet.DeviceMappingGoogleChromeProfiles, he.Source)
		case "mdm.user.mdm.idp@example.com":
			require.Equal(t, fleet.DeviceMappingMDMIdpAccounts, he.Source)
		case "mdm.user.custom.installer@example.com":
			require.Equal(t, fleet.DeviceMappingCustomInstaller, he.Source)
		case "mdm.user.custom.override@example.com":
			require.Equal(t, fleet.DeviceMappingCustomOverride, he.Source)
		case "mdm.user.custom.replacement@example.com":
			require.Equal(t, fleet.DeviceMappingCustomReplacement, he.Source)
		default:
			t.Fatalf("unexpected email '%s'", he.Email)
		}
	}
}
