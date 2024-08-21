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

package release

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gotofu.com/mochi/change"
	"gotofu.com/mochi/config"
	"gotofu.com/mochi/domain"
	"gotofu.com/mochi/utils/git"
)

func Get(target *domain.Target) []*domain.ReleaseNote {
	releaseNotes := []*domain.ReleaseNote{}
	releaseNotesByType := make(map[string][]*domain.ReleaseChange)

	files, _ := filepath.Glob(fmt.Sprintf(".mochi/*-%s-*.md", target.Id))
	slog.Debug("Release notes files found.", "files", files)

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			slog.Warn("Error reading change file.", "file", file, "error", err)
			continue
		}

		change, err := change.Parse(string(data))
		if err != nil {
			slog.Warn("Error parsing change file.", "file", file, "error", err)
			continue
		}

		releaseNotesByType[change.Type.Id] = append(releaseNotesByType[change.Type.Id], &domain.ReleaseChange{
			Change: change,
			File:   file,
		})
	}

	for _, t := range config.Configuration.Types {
		releaseNotesForType := releaseNotesByType[t.Id]
		slog.Debug("Release notes found for type.", "type", t.Id, "count", len(releaseNotesForType))
		if len(releaseNotesForType) > 0 {
			releaseNotes = append(releaseNotes, &domain.ReleaseNote{
				Type:    &t,
				Changes: releaseNotesForType,
			})
		}
	}

	return releaseNotes
}

func Commit(release *domain.Release, rebase bool) error {
	if err := git.EnsureClean(); err != nil {
		return err
	}

	for _, note := range release.Notes {
		for _, change := range note.Changes {
			slog.Debug("Cleaning up release note file.", "file", change.File)

			if err := os.Remove(change.File); err != nil {
				return err
			}

			if err := git.Add(change.File); err != nil {
				return err
			}
		}
	}

	if err := git.Commit(fmt.Sprintf("chore: release %s", release.Tag.String())); err != nil {
		return err
	}

	if err := git.Tag(release.Tag.String()); err != nil {
		return err
	}

	if err := git.Checkout(config.Configuration.BaseBranch, ""); err != nil {
		return err
	}

	if rebase {
		if err := git.Rebase(release.Tag.Branch()); err != nil {
			return err
		}
	} else {
		if err := git.Merge(release.Tag.Branch()); err != nil {
			return err
		}
	}

	if err := git.DeleteBranch(release.Tag.Branch()); err != nil {
		return err
	}

	return nil
}
