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

package domain

import (
	"fmt"
	"io"
	"text/template"
	"time"
)

var changeTemplate, _ = template.New("change").Parse(
	`---
target: {{ .Target.Id }}
type: {{ .Type.Id }}
---

{{ .Message }}
`)

type Change struct {
	Type    *ChangeType
	Target  *Target
	Message string
}

func (c Change) Filename() string {
	return fmt.Sprintf("%s-%s-%s.md", time.Now().Format("20060102150405"), c.Target.Id, c.Type.Id)
}

func (c Change) Render(wr io.Writer) error {
	return changeTemplate.Execute(wr, c)
}
