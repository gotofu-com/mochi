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

func TestTarget_GetTagPrefix(t *testing.T) {
	customTagPrefix := "custom-api"

	tests := []struct {
		name     string
		target   Target
		expected string
	}{
		{
			name: "returns custom tagPrefix when set",
			target: Target{
				Name:      "API",
				Id:        "api",
				TagPrefix: &customTagPrefix,
			},
			expected: "custom-api",
		},
		{
			name: "returns Id when tagPrefix is nil",
			target: Target{
				Name:      "API",
				Id:        "api",
				TagPrefix: nil,
			},
			expected: "api",
		},
		{
			name: "returns Id when tagPrefix is not set",
			target: Target{
				Name: "Frontend",
				Id:   "frontend",
			},
			expected: "frontend",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.target.GetTagPrefix()
			if result != tt.expected {
				t.Errorf("GetTagPrefix() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestTarget_GetTicketPrefix(t *testing.T) {
	customTicketPrefix := "APIPROJ"
	emptyTicketPrefix := ""

	tests := []struct {
		name     string
		target   Target
		expected string
	}{
		{
			name: "returns custom ticketPrefix when set",
			target: Target{
				Name:         "API Service",
				Id:           "api",
				TicketPrefix: &customTicketPrefix,
			},
			expected: "APIPROJ",
		},
		{
			name: "returns Name when ticketPrefix is nil",
			target: Target{
				Name:         "API Service",
				Id:           "api",
				TicketPrefix: nil,
			},
			expected: "API Service",
		},
		{
			name: "returns Name when ticketPrefix is not set",
			target: Target{
				Name: "Frontend App",
				Id:   "frontend",
			},
			expected: "Frontend App",
		},
		{
			name: "handles empty custom ticketPrefix",
			target: Target{
				Name:         "Backend",
				Id:           "backend",
				TicketPrefix: &emptyTicketPrefix,
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.target.GetTicketPrefix()
			if result != tt.expected {
				t.Errorf("GetTicketPrefix() = %q, want %q", result, tt.expected)
			}
		})
	}
}
