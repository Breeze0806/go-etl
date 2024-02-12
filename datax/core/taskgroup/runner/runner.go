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

package runner

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

// Runner - A component responsible for executing tasks or jobs
type Runner interface {
	Plugin() plugin.Task           // Plugin Task - A task associated with a plugin
	Shutdown() error               // Close - Shuts down or terminates the operation of the runner
	Run(ctx context.Context) error // Run - Initiates the execution of a task or job by the runner
}

type baseRunner struct {
}
