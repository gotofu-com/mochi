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
	"time"
)

type Version struct {
	Year  int
	Week  int
	Patch int
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Year, v.Week, v.Patch)
}

func (v Version) Branch(t *Target) string {
	return fmt.Sprintf("release/%s/%s", t.Id, v.String())
}

func (v Version) Validate() error {
	if v.Year <= 0 || v.Year > time.Now().Year() {
		return fmt.Errorf("year is not in the expected range")
	}
	if v.Week <= 0 || v.Week > 53 {
		return fmt.Errorf("week is not in the expected range")
	}
	if v.Patch < 0 {
		return fmt.Errorf("patch is not in the expected range")
	}
	return nil
}

func (v Version) Bump() *Version {
	return &Version{
		Year:  v.Year,
		Week:  v.Week,
		Patch: v.Patch + 1,
	}
}

func (v Version) IsSameWeek(other *Version) bool {
	return v.Year == other.Year && v.Week == other.Week
}
