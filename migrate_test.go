// Copyright (C) 2023  Allen Li
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package migrate

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Example() {
	s := NewMigrationSet([]Migration{
		{
			From: 0,
			To:   1,
			Func: func(ctx context.Context, d *sql.DB) error {
				t, err := d.Begin()
				if err != nil {
					return err
				}
				defer t.Rollback()
				_, err = t.Exec(`
CREATE TABLE user (
    id INTEGER,
    name TEXT,
    PRIMARY KEY (id)
)`)
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			From: 1,
			To:   2,
			Func: func(ctx context.Context, d *sql.DB) error {
				t, err := d.Begin()
				if err != nil {
					return err
				}
				defer t.Rollback()
				_, err = t.Exec(`ALTER TABLE user ADD COLUMN description TEXT`)
				if err != nil {
					return err
				}
				return nil
			},
		},
	})

	d, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		panic(err)
	}
	_ = s.Migrate(context.Background(), d)
}

func TestEmptyMigration(t *testing.T) {
	d := testDB(t)
	defer d.Close()
	s := NewMigrationSet(nil)
	if err := s.Migrate(context.Background(), d); err != nil {
		t.Errorf("Error migrating database: %s", err)
	}
}

func TestSimpleMigration(t *testing.T) {
	d := testDB(t)
	defer d.Close()
	var (
		s0 spyFunc
	)
	s := NewMigrationSet([]Migration{{
		From: 0,
		To:   1,
		Func: s0.migrate,
	}})
	if err := s.Migrate(context.Background(), d); err != nil {
		t.Errorf("Error migrating database: %s", err)
	}
	if !s0.called {
		t.Errorf("Migration function not called")
	}
}

func testDB(t *testing.T) *sql.DB {
	// Cannot be used concurrently!
	d, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("Error opening database: %s", err)
	}
	d.SetConnMaxLifetime(0)
	return d
}

type spyFunc struct {
	called bool
}

func (s *spyFunc) migrate(ctx context.Context, d *sql.DB) error {
	s.called = true
	return nil
}
