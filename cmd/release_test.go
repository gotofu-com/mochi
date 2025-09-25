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
	"os"
	"path/filepath"
	"testing"

	"gotofu.com/mochi/config"
	"gotofu.com/mochi/domain"
	"gotofu.com/mochi/release"

	"github.com/spf13/viper"
)

func TestReleaseCommandsWorkInSubfolders(t *testing.T) {
	// Create temporary directory structure
	tempDir, err := os.MkdirTemp("", "mochi-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .mochi directory and config
	mochiDir := filepath.Join(tempDir, ".mochi")
	if err := os.MkdirAll(mochiDir, 0755); err != nil {
		t.Fatalf("Failed to create .mochi dir: %v", err)
	}

	configContent := `baseBranch: main
types:
  - id: feature
    name: Feature
    title: Features
  - id: bugfix
    name: Bug fix
    title: Bug Fixes
targets:
  - id: backend
    name: Backend Service
    tagPrefix: backend
    ticketPrefix: BE
`

	configFile := filepath.Join(mochiDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Create some sample change files
	changeFiles := []struct {
		filename string
		content  string
	}{
		{
			filename: "20240101120000-backend-feature.md",
			content: `---
target: backend
type: feature
ticketId: "123"
---
Add new authentication endpoint`,
		},
		{
			filename: "20240101130000-backend-bugfix.md",
			content: `---
target: backend
type: bugfix
ticketId: "456"
---
Fix memory leak in user service`,
		},
	}

	for _, cf := range changeFiles {
		changeFilePath := filepath.Join(mochiDir, cf.filename)
		if err := os.WriteFile(changeFilePath, []byte(cf.content), 0644); err != nil {
			t.Fatalf("Failed to write change file %s: %v", cf.filename, err)
		}
	}

	// Create subfolder structure
	subDir := filepath.Join(tempDir, "some", "deep", "subfolder")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subfolder: %v", err)
	}

	// Save current working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	tests := []struct {
		name        string
		workingDir  string
		expectFiles int
	}{
		{
			name:        "release works from root directory",
			workingDir:  tempDir,
			expectFiles: 2,
		},
		{
			name:        "release works from subfolder",
			workingDir:  subDir,
			expectFiles: 2,
		},
		{
			name:        "release works from intermediate folder",
			workingDir:  filepath.Join(tempDir, "some"),
			expectFiles: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Change to test directory
			if err := os.Chdir(tt.workingDir); err != nil {
				t.Fatalf("Failed to change to test directory: %v", err)
			}

			// Reset viper configuration
			viper.Reset()

			// Initialize config (this will discover the config directory)
			config.InitConfig()

			// Verify configuration was loaded
			if config.Configuration == nil {
				t.Fatal("Configuration not loaded")
			}

			if len(config.Configuration.Targets) == 0 {
				t.Fatal("No targets found in configuration")
			}

			// Find the backend target
			var backendTarget *domain.Target
			for i, target := range config.Configuration.Targets {
				if target.Id == "backend" {
					backendTarget = &config.Configuration.Targets[i]
					break
				}
			}

			if backendTarget == nil {
				t.Fatal("Backend target not found in configuration")
			}

			// Test that release.Get finds the change files
			releaseNotes := release.Get(backendTarget)

			if len(releaseNotes) == 0 {
				t.Errorf("Expected release notes to be found, got none")
			}

			// Count total changes across all release notes
			totalChanges := 0
			for _, note := range releaseNotes {
				totalChanges += len(note.Changes)
			}

			if totalChanges != tt.expectFiles {
				t.Errorf("Expected %d changes, got %d", tt.expectFiles, totalChanges)
			}

			// Verify we found both change types
			foundTypes := make(map[string]bool)
			for _, note := range releaseNotes {
				foundTypes[note.Type.Id] = true
			}

			expectedTypes := []string{"feature", "bugfix"}
			for _, expectedType := range expectedTypes {
				if !foundTypes[expectedType] {
					t.Errorf("Expected to find change type %s", expectedType)
				}
			}
		})
	}

	// Restore original working directory
	if err := os.Chdir(originalWd); err != nil {
		t.Fatalf("Failed to restore original working directory: %v", err)
	}
}

func TestReleasePreviewWorksInSubfolders(t *testing.T) {
	// Create temporary directory structure
	tempDir, err := os.MkdirTemp("", "mochi-preview-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .mochi directory and config
	mochiDir := filepath.Join(tempDir, ".mochi")
	if err := os.MkdirAll(mochiDir, 0755); err != nil {
		t.Fatalf("Failed to create .mochi dir: %v", err)
	}

	configContent := `baseBranch: main
baseTicketUrl: "https://jira.example.com/"
types:
  - id: feature
    name: Feature
    title: Features
targets:
  - id: api
    name: API Service
    tagPrefix: api
    ticketPrefix: API
`

	configFile := filepath.Join(mochiDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Create a sample change file
	changeContent := `---
target: api
type: feature
ticketId: "789"
---
Add user management API`

	changeFile := filepath.Join(mochiDir, "20240101140000-api-feature.md")
	if err := os.WriteFile(changeFile, []byte(changeContent), 0644); err != nil {
		t.Fatalf("Failed to write change file: %v", err)
	}

	// Create deep subfolder
	subDir := filepath.Join(tempDir, "services", "api", "handlers")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subfolder: %v", err)
	}

	// Save current working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Errorf("Failed to restore original working directory: %v", err)
		}
	}()

	// Test from subfolder
	if err := os.Chdir(subDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	// Reset and initialize config
	viper.Reset()
	config.InitConfig()

	// Verify configuration was loaded
	if config.Configuration == nil {
		t.Fatal("Configuration not loaded")
	}

	// Find the API target
	var apiTarget *domain.Target
	for i, target := range config.Configuration.Targets {
		if target.Id == "api" {
			apiTarget = &config.Configuration.Targets[i]
			break
		}
	}

	if apiTarget == nil {
		t.Fatal("API target not found in configuration")
	}

	// Test that release.Get finds the change file from subfolder
	releaseNotes := release.Get(apiTarget)

	if len(releaseNotes) == 0 {
		t.Error("Expected release notes to be found from subfolder, got none")
	}

	// Verify the content
	if len(releaseNotes) > 0 {
		firstNote := releaseNotes[0]
		if firstNote.Type.Id != "feature" {
			t.Errorf("Expected feature type, got %s", firstNote.Type.Id)
		}

		if len(firstNote.Changes) == 0 {
			t.Error("Expected at least one change")
		} else {
			change := firstNote.Changes[0]
			if change.Change.Message != "Add user management API" {
				t.Errorf("Expected 'Add user management API', got '%s'", change.Change.Message)
			}
			if change.Change.TicketId == nil || *change.Change.TicketId != "789" {
				t.Errorf("Expected ticket ID '789', got %v", change.Change.TicketId)
			}
		}
	}
}