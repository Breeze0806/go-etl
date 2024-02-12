// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package database

import (
	"fmt"
	"sync"
)

// Dialect Database Dialect
type Dialect interface {
	Source(*BaseSource) (Source, error) // Data Source
}

var dialects = &dialectMap{
	dialects: make(map[string]Dialect),
}

// RegisterDialect Registers a database dialect. A panic occurs when the registered name is the same or the dialect is empty.
func RegisterDialect(name string, dialect Dialect) {
	if err := dialects.register(name, dialect); err != nil {
		panic(err)
	}
}

// UnregisterAllDialects Unregisters all database dialects.
func UnregisterAllDialects() {
	dialects.unregisterAll()
}

type dialectMap struct {
	sync.RWMutex
	dialects map[string]Dialect
}

func (d *dialectMap) register(name string, dialect Dialect) error {
	if dialect == nil {
		return fmt.Errorf("dialect %v is nil", name)
	}

	d.Lock()
	defer d.Unlock()
	if _, ok := d.dialects[name]; ok {
		return fmt.Errorf("dialect %v exists", name)
	}

	d.dialects[name] = dialect
	return nil
}

func (d *dialectMap) dialect(name string) (dialect Dialect, ok bool) {
	d.RLock()
	defer d.RUnlock()
	dialect, ok = d.dialects[name]
	return
}

func (d *dialectMap) unregisterAll() {
	d.Lock()
	defer d.Unlock()
	d.dialects = make(map[string]Dialect)
}
