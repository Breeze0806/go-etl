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

package db2

import "golang.org/x/text/encoding/simplifiedchinese"

//db2的char和varchar类型在windows下中文字符集是gbk
func decodeChinese(data []byte) ([]byte, error) {
	return simplifiedchinese.GBK.NewDecoder().Bytes(data)
}
