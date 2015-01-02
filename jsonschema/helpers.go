// Copyright 2013 sigu-399 ( https://github.com/sigu-399 )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// author           sigu-399
// author-github    https://github.com/sigu-399
// author-mail      sigu.399@gmail.com
//
// repository-name  jsonschema
// repository-desc  An implementation of JSON Schema, based on IETF's draft v4 - Go language.
//
// description		Helper functions to load JSON files.
//
// created          26-02-2013

package jsonschema

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// GetHTTPJSON reads a JSON document from a HTTP request
// HTTP request must return 200 OK
func GetHTTPJSON(url string) (interface{}, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Could not access schema " + resp.Status)
	}

	bodyBuff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var document interface{}
	err = json.Unmarshal(bodyBuff, &document)
	if err != nil {
		return nil, err
	}

	return document, nil
}

// GetFileJSON reads a JSON document from a file path
func GetFileJSON(filepath string) (interface{}, error) {

	bodyBuff, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var document interface{}
	err = json.Unmarshal(bodyBuff, &document)
	if err != nil {
		return nil, err
	}

	return document, nil
}
