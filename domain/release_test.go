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
	"testing"
)

func TestReleaseTemplate_TicketPrefixFormatting(t *testing.T) {
	customTicketPrefix := "APIPROJ"
	baseTicketUrl := "https://jira.example.com/"
	ticketId := "123"
	
	tests := []struct {
		name           string
		release        Release
		expectedOutput string
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
			expectedOutput: `
## Features
- Add new authentication endpoint - [APIPROJ-123](https://jira.example.com/APIPROJ-123)
`,
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
			expectedOutput: `
## Bug Fixes
- Fix login form validation - [Frontend App-456](https://github.com/example/repo/issues/Frontend App-456)
`,
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
			expectedOutput: `
## Features
- Improve database performance
`,
		},
		{
			name: "handles multiple changes in same category",
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
										Name:         "API",
										Id:           "api",
										TicketPrefix: &customTicketPrefix,
									},
									Message:  "Add user profile endpoint",
									TicketId: &[]string{"101"}[0],
								},
							},
							{
								Change: &Change{
									Type: &ChangeType{
										Id:    "feature",
										Name:  "Feature",
										Title: "Features",
									},
									Target: &Target{
										Name: "API",
										Id:   "api",
										TicketPrefix: &customTicketPrefix,
									},
									Message:  "Add settings endpoint",
									TicketId: &[]string{"102"}[0],
								},
							},
						},
					},
				},
			},
			expectedOutput: `
## Features
- Add user profile endpoint - [APIPROJ-101](https://jira.example.com/APIPROJ-101)
- Add settings endpoint - [APIPROJ-102](https://jira.example.com/APIPROJ-102)
`,
		},
		{
			name: "handles multiple note categories",
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
										Name: "Backend",
										Id:   "backend",
									},
									Message:  "Add caching layer",
									TicketId: nil,
								},
							},
						},
					},
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
										Name: "Backend",
										Id:   "backend",
									},
									Message:  "Fix memory leak",
									TicketId: &[]string{"999"}[0],
								},
							},
						},
					},
				},
			},
			expectedOutput: `
## Features
- Add caching layer

## Bug Fixes
- Fix memory leak - [Backend-999](https://jira.example.com/Backend-999)
`,
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
			if result != tt.expectedOutput {
				t.Errorf("Release output mismatch.\nExpected:\n%q\n\nGot:\n%q", tt.expectedOutput, result)
			}
		})
	}
}