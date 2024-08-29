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

package version

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"gotofu.com/mochi/domain"
	"gotofu.com/mochi/utils/git"
)

var versionRegex = regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)

func Parse(version string) (*domain.Version, error) {
	var (
		v   domain.Version
		err error
	)

	parts := versionRegex.FindStringSubmatch(version)
	if len(parts) != 4 {
		return nil, fmt.Errorf("version is not in the format year.week.patch")
	}

	if v.Year, err = strconv.Atoi(parts[1]); err != nil {
		return nil, fmt.Errorf("year is not a valid number")
	}
	if v.Week, err = strconv.Atoi(parts[2]); err != nil {
		return nil, fmt.Errorf("week is not a valid number")
	}
	if v.Patch, err = strconv.Atoi(parts[3]); err != nil {
		return nil, fmt.Errorf("patch is not a valid number")
	}

	if err = v.Validate(); err != nil {
		return nil, err
	}

	return &v, nil
}

func Next(target *domain.Target, latestVersion *domain.Version) *domain.Version {
	currentYear, currentWeek := time.Now().ISOWeek()
	nextVersion := &domain.Version{
		Year:  currentYear,
		Week:  currentWeek,
		Patch: 0,
	}

	if latestVersion != nil && latestVersion.IsSameWeek(nextVersion) {
		nextVersion = latestVersion.Bump()
	}

	return nextVersion
}

func Latest(target *domain.Target) (*domain.Version, error) {
	latestGitTag, err := git.LatestTagForTarget(target.Id)
	if err != nil {
		return nil, err
	}

	latestVersion, err := Parse(string(latestGitTag))
	if err != nil {
		return nil, err
	}

	return latestVersion, nil
}
