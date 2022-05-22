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

package csv

import (
	"golang.org/x/text/encoding/simplifiedchinese"
)

var (
	encoders = map[string]encode{
		"gbk":   gbkEncoder,
		"utf-8": utf8Encoder,
	}
	decoders = map[string]decode{
		"gbk":   gbkDecoder,
		"utf-8": utf8Decoder,
	}
)

type encode func(string) (string, error)

type decode func(string) (string, error)

func gbkDecoder(src string) (dest string, err error) {
	dest, err = simplifiedchinese.GBK.NewDecoder().String(src)
	return
}

func utf8Decoder(src string) (dest string, err error) {
	return src, nil
}

func gbkEncoder(src string) (dest string, err error) {
	dest, err = simplifiedchinese.GBK.NewEncoder().String(src)
	return
}

func utf8Encoder(src string) (dest string, err error) {
	return src, nil
}
