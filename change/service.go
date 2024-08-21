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

package change

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gotofu.com/mochi/change_type"
	"gotofu.com/mochi/domain"
	"gotofu.com/mochi/target"

	"gopkg.in/yaml.v3"
)

func Commit(c *domain.Change) error {
	fileName := fmt.Sprintf(".mochi/%s", c.Filename())
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	return c.Render(file)
}

type ChangeMeta struct {
	Target string
	Type   string
}

var regex = regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)$`)

func Parse(rawChange string) (*domain.Change, error) {
	var (
		c   domain.Change
		m   ChangeMeta
		err error
	)

	matches := regex.FindStringSubmatch(rawChange)
	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid change format")
	}

	if err := yaml.Unmarshal([]byte(matches[1]), &m); err != nil {
		return nil, err
	}

	c.Message = strings.TrimSpace(matches[2])
	if c.Target, err = target.Get(m.Target); err != nil {
		return nil, err
	}
	if c.Type, err = change_type.Get(m.Type); err != nil {
		return nil, err
	}

	return &c, nil
}
