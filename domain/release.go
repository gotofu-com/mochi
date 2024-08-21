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
	"io"
	"text/template"
)

type ReleaseChange struct {
	Change *Change
	File   string
}

type ReleaseNote struct {
	Type    *ChangeType
	Changes []*ReleaseChange
}

type Release struct {
	Tag   *Tag
	Notes []*ReleaseNote
}

var releaseTemplate = template.Must(template.New("release").Parse(`
{{- range .Notes }}
## {{ .Type.Title }}
{{ range .Changes -}}
- {{ .Change.Message }}
{{ end -}}
{{ end -}}
`))

func (r Release) Render(wr io.Writer) error {
	return releaseTemplate.Execute(wr, r)
}
