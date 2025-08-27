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

import "testing"

func TestTag_String_FormattingWithTagPrefix(t *testing.T) {
	customTagPrefix := "custom-api"
	complexTagPrefix := "mobile-ios-app"

	tests := []struct {
		name           string
		tag            Tag
		expectedString string
	}{
		{
			name: "uses custom tagPrefix in tag string",
			tag: Tag{
				Target: &Target{
					Name:      "API Service",
					Id:        "api",
					TagPrefix: &customTagPrefix,
				},
				Version: &Version{
					Year:  2025,
					Week:  12,
					Patch: 3,
				},
			},
			expectedString: "custom-api@2025.12.3",
		},
		{
			name: "uses Id as tagPrefix when not set",
			tag: Tag{
				Target: &Target{
					Name: "Frontend",
					Id:   "frontend",
				},
				Version: &Version{
					Year:  2025,
					Week:  8,
					Patch: 1,
				},
			},
			expectedString: "frontend@2025.8.1",
		},
		{
			name: "uses Id as tagPrefix when tagPrefix is nil",
			tag: Tag{
				Target: &Target{
					Name:      "Backend Service",
					Id:        "backend",
					TagPrefix: nil,
				},
				Version: &Version{
					Year:  2024,
					Week:  52,
					Patch: 10,
				},
			},
			expectedString: "backend@2024.52.10",
		},
		{
			name: "handles complex custom tagPrefix",
			tag: Tag{
				Target: &Target{
					Name:      "Mobile App",
					Id:        "mobile",
					TagPrefix: &complexTagPrefix,
				},
				Version: &Version{
					Year:  2025,
					Week:  1,
					Patch: 0,
				},
			},
			expectedString: "mobile-ios-app@2025.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tag.String()
			if result != tt.expectedString {
				t.Errorf("Tag.String() = %q, want %q", result, tt.expectedString)
			}
		})
	}
}
