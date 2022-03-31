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

//  DBStorage（数据库存储）
//     +--------------------------------------------------------+
//     |                                                        |
//   MasterTable（主表）---------+-------+----------------- SlaveTable（从表）
//     |          	            |       |
//     |                        |       +--compare(比较)--+
//     |                        |      /                  |
//     |                        |     /                   +---TableDiffer(表差异)
//     |				    TableNameMap（表映射关联）                |
//   Tracker（跟踪器）                                         read/writer(读/写)
//     +------------------+                                          |
//     |                  |                                 DifferStorage（差异存储）
// OffsetTracker   PageParamTracker
//  (位移记录器)      (分页记录器)
package plugin
