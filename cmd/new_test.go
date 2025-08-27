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

package cmd

import (
	"testing"
	"gotofu.com/mochi/domain"
)

func TestGetTargetRegex_FormattingWithTicketPrefix(t *testing.T) {
	customTicketPrefix := "APIPROJ"
	anotherTicketPrefix := "FRONTEND"
	
	tests := []struct {
		name           string
		targets        []domain.Target
		expectedRegex  string
	}{
		{
			name: "generates regex using custom ticketPrefix",
			targets: []domain.Target{
				{
					Name:         "API Service",
					Id:           "api",
					TicketPrefix: &customTicketPrefix,
				},
			},
			expectedRegex: `(?i)(?:APIPROJ)-(\d+)`,
		},
		{
			name: "generates regex using Name when ticketPrefix not set",
			targets: []domain.Target{
				{
					Name: "Frontend App",
					Id:   "frontend",
				},
			},
			expectedRegex: `(?i)(?:Frontend App)-(\d+)`,
		},
		{
			name: "generates regex for multiple targets with mixed prefixes",
			targets: []domain.Target{
				{
					Name:         "API Service",
					Id:           "api",
					TicketPrefix: &customTicketPrefix,
				},
				{
					Name:         "Frontend App",
					Id:           "frontend",
					TicketPrefix: &anotherTicketPrefix,
				},
				{
					Name: "Backend Service",
					Id:   "backend",
				},
			},
			expectedRegex: `(?i)(?:APIPROJ|FRONTEND|Backend Service)-(\d+)`,
		},
		{
			name: "handles empty targets slice",
			targets: []domain.Target{},
			expectedRegex: `(?i)(?:)-(\d+)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTargetRegex(tt.targets)
			if result != tt.expectedRegex {
				t.Errorf("getTargetRegex() = %q, want %q", result, tt.expectedRegex)
			}
		})
	}
}