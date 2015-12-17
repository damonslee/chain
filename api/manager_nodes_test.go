package api

import (
	"testing"

	"chain/api/asset"
	"chain/database/pg"
	"chain/database/pg/pgtest"
	"chain/errors"
	"chain/net/http/authn"
)

func TestCreateManagerNodeBadXPub(t *testing.T) {
	ctx := pgtest.NewContext(t, testUserFixture, `
		INSERT INTO projects(id, name) VALUES ('a1', 'x');
		INSERT INTO members (project_id, user_id, role)
			VALUES ('a1', 'sample-user-id-0', 'admin');
	`)
	defer pgtest.Finish(ctx)
	ctx = authn.NewContext(ctx, "sample-user-id-0")

	req := map[string]interface{}{
		"label":               "node",
		"signatures_required": 1,
		"keys":                []*asset.CreateNodeKeySpec{{Type: "node", XPub: "badxpub"}},
	}

	_, err := createManagerNode(ctx, "a1", req)
	if got := errors.Root(err); got != asset.ErrBadKeySpec {
		t.Fatalf("err = %v want %v", got, asset.ErrBadKeySpec)
	}
}

func TestCreateManagerNode(t *testing.T) {
	ctx := pgtest.NewContext(t, testUserFixture, `
		INSERT INTO projects(id, name) VALUES ('a1', 'x');
		INSERT INTO members (project_id, user_id, role)
			VALUES ('a1', 'sample-user-id-0', 'admin');
	`)
	defer pgtest.Finish(ctx)
	ctx = authn.NewContext(ctx, "sample-user-id-0")

	req := map[string]interface{}{
		"label":               "node",
		"keys":                []*asset.CreateNodeKeySpec{{Type: "node", Generate: true}},
		"signatures_required": 1,
	}

	_, err := createManagerNode(ctx, "a1", req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	var count int
	var checkQ = `SELECT COUNT(*) FROM manager_nodes`
	err = pg.FromContext(ctx).QueryRow(ctx, checkQ).Scan(&count)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if count != 1 {
		t.Fatalf("checking results, want=1, got=%d", count)
	}
}

func TestCreateManagerNodeDeprecated(t *testing.T) {
	ctx := pgtest.NewContext(t, testUserFixture, `
		INSERT INTO projects(id, name) VALUES ('a1', 'x');
		INSERT INTO members (project_id, user_id, role)
			VALUES ('a1', 'sample-user-id-0', 'admin');
	`)
	defer pgtest.Finish(ctx)
	ctx = authn.NewContext(ctx, "sample-user-id-0")

	req := map[string]interface{}{
		"label":        "deprecated node",
		"generate_key": true,
	}

	_, err := createManagerNode(ctx, "a1", req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	var count int
	var checkQ = `SELECT COUNT(*) FROM manager_nodes`
	err = pg.FromContext(ctx).QueryRow(ctx, checkQ).Scan(&count)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if count != 1 {
		t.Fatalf("checking results, want=1, got=%d", count)
	}
}
