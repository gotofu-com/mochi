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
	"testing"
	"time"

	"gotofu.com/mochi/domain"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    *domain.Version
		expectError bool
	}{
		{
			name:  "parses valid version",
			input: "2025.44.0",
			expected: &domain.Version{
				Year:  2025,
				Week:  44,
				Patch: 0,
			},
			expectError: false,
		},
		{
			name:  "parses version with patch number",
			input: "2025.44.5",
			expected: &domain.Version{
				Year:  2025,
				Week:  44,
				Patch: 5,
			},
			expectError: false,
		},
		{
			name:  "parses version string only (no prefix)",
			input: "2025.44.0",
			expected: &domain.Version{
				Year:  2025,
				Week:  44,
				Patch: 0,
			},
			expectError: false,
		},
		{
			name:        "fails on invalid format",
			input:       "invalid",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "fails on incomplete version",
			input:       "2025.44",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "fails on invalid week number",
			input:       "2025.99.0",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "fails on negative patch",
			input:       "2025.44.-1",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Parse(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
				return
			}

			if result.Year != tt.expected.Year || result.Week != tt.expected.Week || result.Patch != tt.expected.Patch {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNext(t *testing.T) {
	currentYear, currentWeek := time.Now().ISOWeek()

	tests := []struct {
		name          string
		target        *domain.Target
		latestVersion *domain.Version
		expected      *domain.Version
	}{
		{
			name: "creates first version when no latest version exists",
			target: &domain.Target{
				Id:   "backend",
				Name: "Backend Service",
			},
			latestVersion: nil,
			expected: &domain.Version{
				Year:  currentYear,
				Week:  currentWeek,
				Patch: 0,
			},
		},
		{
			name: "bumps patch when latest version is same week",
			target: &domain.Target{
				Id:   "backend",
				Name: "Backend Service",
			},
			latestVersion: &domain.Version{
				Year:  currentYear,
				Week:  currentWeek,
				Patch: 0,
			},
			expected: &domain.Version{
				Year:  currentYear,
				Week:  currentWeek,
				Patch: 1,
			},
		},
		{
			name: "bumps patch multiple times in same week",
			target: &domain.Target{
				Id:   "backend",
				Name: "Backend Service",
			},
			latestVersion: &domain.Version{
				Year:  currentYear,
				Week:  currentWeek,
				Patch: 3,
			},
			expected: &domain.Version{
				Year:  currentYear,
				Week:  currentWeek,
				Patch: 4,
			},
		},
		{
			name: "resets patch when moving to new week",
			target: &domain.Target{
				Id:   "backend",
				Name: "Backend Service",
			},
			latestVersion: &domain.Version{
				Year:  currentYear,
				Week:  currentWeek - 1,
				Patch: 5,
			},
			expected: &domain.Version{
				Year:  currentYear,
				Week:  currentWeek,
				Patch: 0,
			},
		},
		{
			name: "resets patch when moving to new year",
			target: &domain.Target{
				Id:   "backend",
				Name: "Backend Service",
			},
			latestVersion: &domain.Version{
				Year:  currentYear - 1,
				Week:  52,
				Patch: 10,
			},
			expected: &domain.Version{
				Year:  currentYear,
				Week:  currentWeek,
				Patch: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Next(tt.target, tt.latestVersion)

			if result.Year != tt.expected.Year || result.Week != tt.expected.Week || result.Patch != tt.expected.Patch {
				t.Errorf("Next() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNext_SameWeekBugScenario(t *testing.T) {
	// This test specifically addresses the bug where multiple releases
	// in the same week would fail because the patch number wasn't incremented.
	//
	// Note: This test validates the Next() logic. The Latest() function that
	// retrieves existing tags must use target.GetTagPrefix() (not target.Id)
	// to ensure tags are found using the same prefix they were created with.
	currentYear, currentWeek := time.Now().ISOWeek()

	target := &domain.Target{
		Id:   "bpo",
		Name: "Backend Platform",
	}

	// First release of the week - should be patch 0
	firstRelease := Next(target, nil)
	if firstRelease.Patch != 0 {
		t.Errorf("First release patch = %d, want 0", firstRelease.Patch)
	}

	// Second release of the same week - should bump to patch 1
	secondRelease := Next(target, firstRelease)
	if secondRelease.Year != currentYear || secondRelease.Week != currentWeek {
		t.Errorf("Second release week = %d.%d, want %d.%d",
			secondRelease.Year, secondRelease.Week, currentYear, currentWeek)
	}
	if secondRelease.Patch != 1 {
		t.Errorf("Second release patch = %d, want 1", secondRelease.Patch)
	}

	// Third release of the same week - should bump to patch 2
	thirdRelease := Next(target, secondRelease)
	if thirdRelease.Patch != 2 {
		t.Errorf("Third release patch = %d, want 2", thirdRelease.Patch)
	}
}

func TestVersion_IsSameWeek(t *testing.T) {
	tests := []struct {
		name     string
		version1 domain.Version
		version2 *domain.Version
		expected bool
	}{
		{
			name: "same year and week",
			version1: domain.Version{
				Year:  2025,
				Week:  44,
				Patch: 0,
			},
			version2: &domain.Version{
				Year:  2025,
				Week:  44,
				Patch: 5,
			},
			expected: true,
		},
		{
			name: "different week same year",
			version1: domain.Version{
				Year:  2025,
				Week:  44,
				Patch: 0,
			},
			version2: &domain.Version{
				Year:  2025,
				Week:  45,
				Patch: 0,
			},
			expected: false,
		},
		{
			name: "different year same week",
			version1: domain.Version{
				Year:  2024,
				Week:  44,
				Patch: 0,
			},
			version2: &domain.Version{
				Year:  2025,
				Week:  44,
				Patch: 0,
			},
			expected: false,
		},
		{
			name: "completely different",
			version1: domain.Version{
				Year:  2024,
				Week:  1,
				Patch: 0,
			},
			version2: &domain.Version{
				Year:  2025,
				Week:  52,
				Patch: 10,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.version1.IsSameWeek(tt.version2)
			if result != tt.expected {
				t.Errorf("IsSameWeek() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVersion_Bump(t *testing.T) {
	tests := []struct {
		name     string
		version  domain.Version
		expected *domain.Version
	}{
		{
			name: "bumps patch from 0",
			version: domain.Version{
				Year:  2025,
				Week:  44,
				Patch: 0,
			},
			expected: &domain.Version{
				Year:  2025,
				Week:  44,
				Patch: 1,
			},
		},
		{
			name: "bumps patch from higher number",
			version: domain.Version{
				Year:  2025,
				Week:  44,
				Patch: 5,
			},
			expected: &domain.Version{
				Year:  2025,
				Week:  44,
				Patch: 6,
			},
		},
		{
			name: "preserves year and week",
			version: domain.Version{
				Year:  2024,
				Week:  1,
				Patch: 99,
			},
			expected: &domain.Version{
				Year:  2024,
				Week:  1,
				Patch: 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.version.Bump()
			if result.Year != tt.expected.Year || result.Week != tt.expected.Week || result.Patch != tt.expected.Patch {
				t.Errorf("Bump() = %v, want %v", result, tt.expected)
			}
		})
	}
}
