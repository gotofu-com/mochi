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
)

type Tag struct {
	Target  *Target
	Version *Version
}

func (t Tag) String() string {
	return fmt.Sprintf("%s@%d.%d.%d", t.Target.Id, t.Version.Year, t.Version.Week, t.Version.Patch)
}

func (t Tag) Branch() string {
	return t.Version.Branch(t.Target)
}
