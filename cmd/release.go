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
	"fmt"
	"log/slog"
	"os"

	"gotofu.com/mochi/config"
	"gotofu.com/mochi/domain"
	"gotofu.com/mochi/release"
	"gotofu.com/mochi/tag"
	"gotofu.com/mochi/target"
	"gotofu.com/mochi/utils/git"
	"gotofu.com/mochi/version"

	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:   "release [command]",
	Short: "Create and manage releases",
}

var releaseStartCmd = &cobra.Command{
	Use:   "start [target]",
	Short: "Start working on a release",
	Long: `The "start" command creates a new branch following the format below:

release/<target>/<version>`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var comps []string

		switch len(args) {
		case 0:
			for _, t := range config.Configuration.Targets {
				comps = append(comps, t.Id)
			}
		case 1:
			comps = cobra.AppendActiveHelp(comps, "This command does not take any more arguments (but may accept flags)")
		default:
			comps = cobra.AppendActiveHelp(comps, "ERROR: Too many arguments specified")
		}

		return comps, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			currentTarget *domain.Target
			latestVersion *domain.Version
			gitBase       string
			err           error
		)

		if currentTarget, err = target.Get(args[0]); err != nil {
			return err
		}

		if baseFlag, _ := cmd.Flags().GetString("base"); baseFlag != "" {
			if baseVersion, err := version.Parse(baseFlag); err != nil {
				return err
			} else {
				gitBase = domain.Tag{
					Target:  currentTarget,
					Version: baseVersion,
				}.String()
			}
		}

		if err := git.EnsureClean(); err != nil {
			return err
		}

		if currentBranch, err := git.CurrentBranch(); err != nil {
			return err
		} else if currentBranch != config.Configuration.BaseBranch {
			return fmt.Errorf("you must be on the base branch to start a release")
		}

		if latestVersion, err = version.Latest(currentTarget); err != nil {
			slog.Debug("No valid version found in git tags, falling back to the default current version.", "error", err.Error())
		}

		nextVersion := version.Next(currentTarget, latestVersion)

		if err := git.Checkout(nextVersion.Branch(currentTarget), gitBase); err != nil {
			return err
		}

		fmt.Printf(`Started release %s for %s on branch %s.
		
You can add additional commits in preparation for this release if you wish.

To finalize the release, run 'mochi release finish'. If you wish to preview the release notes, run 'mochi release preview'.
`, nextVersion.String(), currentTarget.Name, nextVersion.Branch(currentTarget))

		return nil
	},
}

var releasePreviewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview the release notes for an in-progress release",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentBranch, err := git.CurrentBranch()
		if err != nil {
			return err
		}
		if currentBranch == config.Configuration.BaseBranch {
			return fmt.Errorf("you must be on a release branch to preview the release notes")
		}

		tag, err := tag.ParseFromBranch(currentBranch)
		if err != nil {
			return err
		}

		releaseNotes := release.Get(tag.Target)
		if len(releaseNotes) == 0 {
			fmt.Println("No release notes found.")
			return nil
		}

		rel := domain.Release{
			Tag:   tag,
			Notes: releaseNotes,
		}
		fmt.Printf("Release notes for %s:\n", tag.String())
		rel.Render(os.Stdout)

		return nil
	},
}

var releaseFinishCmd = &cobra.Command{
	Use:   "finish",
	Short: "Finish a release",
	Long: `The "finish" command finalizes a release by following these steps:
	
1. Gathers the release notes from the release notes files
2. Removes the release notes files, commits the changes, and tags the commit
3. Merges the release branch into the base branch`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rebase, err := cmd.Flags().GetBool("rebase")
		if err != nil {
			return err
		}

		currentBranch, err := git.CurrentBranch()
		if err != nil {
			return err
		}
		if currentBranch == config.Configuration.BaseBranch {
			return fmt.Errorf("you must be on a release branch to finish a release")
		}

		tag, err := tag.ParseFromBranch(currentBranch)
		if err != nil {
			return err
		}

		releaseNotes := release.Get(tag.Target)
		if len(releaseNotes) == 0 {
			return fmt.Errorf("no release notes found; add a release note to finish the release")
		}

		rel := domain.Release{
			Tag:   tag,
			Notes: releaseNotes,
		}
		fmt.Printf("Release notes for %s:\n", tag.String())
		rel.Render(os.Stdout)

		if err := release.Commit(&rel, rebase); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	releaseStartCmd.Flags().StringP("base", "b", "", "the base version to start the release from (e.g. 2024.1.0)")

	releaseFinishCmd.Flags().Bool("rebase", false, "rebase the release branch on top of the base branch instead of merging it")

	releaseCmd.AddCommand(releaseStartCmd)
	releaseCmd.AddCommand(releasePreviewCmd)
	releaseCmd.AddCommand(releaseFinishCmd)

	rootCmd.AddCommand(releaseCmd)
}
