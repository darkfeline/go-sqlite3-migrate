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
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func getUserVersion(d *sql.DB) (int, error) {
	r, err := d.Query("PRAGMA user_version")
	if err != nil {
		return 0, fmt.Errorf("get user version: %s", err)
	}
	defer r.Close()
	ok := r.Next()
	if !ok {
		return 0, fmt.Errorf("get user version: %s", r.Err())
	}
	var v int
	if err := r.Scan(&v); err != nil {
		return 0, fmt.Errorf("get user version: %s", err)
	}
	r.Close()
	if err := r.Err(); err != nil {
		return 0, fmt.Errorf("get user version: %s", err)
	}
	return v, nil
}

func setUserVersion(d *sql.DB, v int) error {
	_, err := d.Exec(fmt.Sprintf("PRAGMA user_version=%d", v))
	if err != nil {
		return fmt.Errorf("set user version %d: %s", v, err)
	}
	return nil
}
