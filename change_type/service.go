/*
Copyright © 2024-present The Mochi Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package change_type

import (
	"fmt"

	"gotofu.com/mochi/config"
	"gotofu.com/mochi/domain"
)

func GetIds() []string {
	var ids []string
	for _, c := range config.Configuration.Types {
		ids = append(ids, c.Id)
	}

	return ids
}

func Get(id string) (*domain.ChangeType, error) {
	for _, t := range config.Configuration.Types {
		if t.Id == id {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("change fragment with ID %s not found", id)
}
