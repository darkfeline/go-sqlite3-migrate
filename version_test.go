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
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestUserVersion(t *testing.T) {
	d := testDB(t)
	defer d.Close()
	v, err := getUserVersion(d)
	if err != nil {
		t.Fatalf("Error getting version: %s", err)
	}
	if v != 0 {
		t.Errorf("Expected 0, got %d", v)
	}
	err = setUserVersion(d, 1)
	if err != nil {
		t.Fatalf("Error setting version: %s", err)
	}
	v, err = getUserVersion(d)
	if err != nil {
		t.Fatalf("Error getting version: %s", err)
	}
	if v != 1 {
		t.Errorf("Expected 1, got %d", v)
	}
}
