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
	"bytes"
	"strings"
	"testing"
)

func TestReleaseTemplate_TicketPrefixFormatting(t *testing.T) {
	customTicketPrefix := "APIPROJ"
	baseTicketUrl := "https://jira.example.com/"
	ticketId := "123"
	
	tests := []struct {
		name             string
		release          Release
		expectedContains []string
	}{
		{
			name: "uses custom ticketPrefix in release template",
			release: Release{
				BaseTicketUrl: &baseTicketUrl,
				Notes: []*ReleaseNote{
					{
						Type: &ChangeType{
							Id:    "feature",
							Name:  "Feature",
							Title: "Features",
						},
						Changes: []*ReleaseChange{
							{
								Change: &Change{
									Type: &ChangeType{
										Id:    "feature",
										Name:  "Feature",
										Title: "Features",
									},
									Target: &Target{
										Name:         "API Service",
										Id:           "api",
										TicketPrefix: &customTicketPrefix,
									},
									Message:  "Add new authentication endpoint",
									TicketId: &ticketId,
								},
							},
						},
					},
				},
			},
			expectedContains: []string{
				"[APIPROJ-123]",
				"(https://jira.example.com/APIPROJ-123)",
				"Add new authentication endpoint",
			},
		},
		{
			name: "uses Name as ticketPrefix when not set",
			release: Release{
				BaseTicketUrl: &[]string{"https://github.com/example/repo/issues/"}[0],
				Notes: []*ReleaseNote{
					{
						Type: &ChangeType{
							Id:    "bugfix",
							Name:  "Bug fix",
							Title: "Bug Fixes",
						},
						Changes: []*ReleaseChange{
							{
								Change: &Change{
									Type: &ChangeType{
										Id:    "bugfix",
										Name:  "Bug fix",
										Title: "Bug Fixes",
									},
									Target: &Target{
										Name: "Frontend App",
										Id:   "frontend",
									},
									Message:  "Fix login form validation",
									TicketId: &[]string{"456"}[0],
								},
							},
						},
					},
				},
			},
			expectedContains: []string{
				"[Frontend App-456]",
				"(https://github.com/example/repo/issues/Frontend App-456)",
				"Fix login form validation",
			},
		},
		{
			name: "handles change without ticketId",
			release: Release{
				BaseTicketUrl: &baseTicketUrl,
				Notes: []*ReleaseNote{
					{
						Type: &ChangeType{
							Id:    "feature",
							Name:  "Feature",
							Title: "Features",
						},
						Changes: []*ReleaseChange{
							{
								Change: &Change{
									Type: &ChangeType{
										Id:    "feature",
										Name:  "Feature",
										Title: "Features",
									},
									Target: &Target{
										Name:         "Backend",
										Id:           "backend",
										TicketPrefix: &customTicketPrefix,
									},
									Message:  "Improve database performance",
									TicketId: nil,
								},
							},
						},
					},
				},
			},
			expectedContains: []string{
				"Improve database performance",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tt.release.Render(&buf)
			if err != nil {
				t.Fatalf("Release.Render() failed: %v", err)
			}

			result := buf.String()
			for _, expected := range tt.expectedContains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected release output to contain %q, but it didn't.\nFull output:\n%s", expected, result)
				}
			}
		})
	}
}