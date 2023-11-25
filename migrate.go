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

// Package migrate implements database migrations SQLite3 using the
// [database/sql] interface.  In theory, it should work with any
// SQLite3 driver, but it is tested against
// [github.com/mattn/go-sqlite3].
//
// This package aims to be simple.  Versions are tracked using the
// user_version PRAGMA.  Only migrations to the latest version are
// supported.
package migrate

import (
	"context"
	"database/sql"
	"fmt"
)

// A MigrationSet contains a set of [Migration] rules which together
// define how to migrate various database versions to a latest
// version.
type MigrationSet struct {
	migrations []Migration
	latest     int
}

// NewMigrationSet creates a new [MigrationSet].
//
// The input slice is used in the returned [MigrationSet] and should
// not be retained by the caller.
//
// The input slice should be sorted in increasing order of
// the "From" field.
//
// If there is more than one [Migration] with the same "From" value,
// then one is chosen arbitrarily when upgrading databases with that
// version.
func NewMigrationSet(m []Migration) *MigrationSet {
	latest := 0
	for _, m := range m {
		latest = max(latest, m.To)
	}
	return &MigrationSet{
		migrations: m,
		latest:     latest,
	}
}

// Migrate migrates the database to the latest version.
func (s *MigrationSet) Migrate(ctx context.Context, d *sql.DB) error {
	v, err := getUserVersion(d)
	if err != nil {
		return fmt.Errorf("migrate: %s", err)
	}
	if v == s.latest {
		return nil
	}
	for _, m := range s.migrations {
		if v != m.From {
			continue
		}
		if err := m.Func(ctx, d); err != nil {
			return fmt.Errorf("migrate from %d to %d: %s", m.From, m.To, err)
		}
		if err := setUserVersion(d, m.To); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
		v = m.To
	}
	if v != s.latest {
		return fmt.Errorf("migrate: no migration path from %d to %d",
			s.latest, v)
	}
	return nil
}

// A Migration describes how to migrate a database from one version to another.
// Only migrations from lower versions to higher versions are supported.
type Migration struct {
	From int
	To   int
	// The migration function must be idempotent.  The function is
	// not wrapped in a transaction, so the function is free to
	// use transactions itself however it wants.
	Func func(context.Context, *sql.DB) error
}
