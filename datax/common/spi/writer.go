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

package spi

import "github.com/Breeze0806/go-etl/datax/common/spi/writer"

// Writer
type Writer interface {
	Job() writer.Job   // Acquire the writing job, which generally cannot be null.
	Task() writer.Task // Acquire the writing task, which generally cannot be null.
}
