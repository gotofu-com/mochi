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

package git

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

func execGit(args ...string) (string, error) {
	slog.Debug("Running git", "args", args)

	var stdout, stderr bytes.Buffer

	cmd := exec.Command("git", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		slog.Debug("Git command failed", "stderr", stderr.String())
		return "", fmt.Errorf("could not run git: %s", stderr.String())
	}

	slog.Debug("Git command succeeded", "stdout", stdout.String())

	return stdout.String(), nil
}

func EnsureClean() error {
	if status, err := execGit("status", "--porcelain"); err != nil {
		return fmt.Errorf("could not check git status: %w", err)
	} else if len(status) > 0 {
		return fmt.Errorf("repository is not clean")
	}

	return nil
}

func RevisionExists(rev string) bool {
	if _, err := execGit("rev-parse", "--verify", rev); err != nil {
		return false
	}

	return true
}

func LatestTagForTarget(target string) (string, error) {
	latestGitTag, err := execGit("describe", "--tags", fmt.Sprintf(`--match=%s*`, target), "--abbrev=0", "HEAD")
	if err != nil {
		return "", fmt.Errorf("could not get latest tag for target %s", target)
	}
	return strings.TrimSpace(latestGitTag), nil
}

func CurrentBranch() (string, error) {
	if result, err := execGit("rev-parse", "--abbrev-ref", "HEAD"); err != nil {
		return "", fmt.Errorf("could not get current branch: %w", err)
	} else {
		return strings.TrimSpace(result), nil
	}
}

func Checkout(branch string, base string) error {
	args := []string{"checkout"}

	if RevisionExists(branch) {
		args = append(args, branch)
	} else {
		args = append(args, "-b", branch)

		if base != "" {
			if RevisionExists(base) {
				args = append(args, base)
			} else {
				return fmt.Errorf("base revision %s does not exist", base)
			}
		}
	}

	if _, err := execGit(args...); err != nil {
		return fmt.Errorf("could not checkout branch %s: %w", branch, err)
	}

	return nil
}

func Tag(tag string) error {
	if _, err := execGit("tag", tag, "-m", tag); err != nil {
		return fmt.Errorf("could not tag %s: %w", tag, err)
	}

	return nil
}

func Add(path string) error {
	if _, err := execGit("add", path); err != nil {
		return fmt.Errorf("could not add %s: %w", path, err)
	}

	return nil
}

func Commit(message string) error {
	if _, err := execGit("commit", "-m", message); err != nil {
		return fmt.Errorf("could not commit: %w", err)
	}

	return nil
}

func Push() error {
	if _, err := execGit("push", "--follow-tags"); err != nil {
		return fmt.Errorf("could not push: %w", err)
	}

	return nil
}

func Merge(branch string) error {
	if _, err := execGit("merge", branch); err != nil {
		return fmt.Errorf("could not merge branch %s: %w", branch, err)
	}

	return nil
}

func Rebase(branch string) error {
	if _, err := execGit("rebase", branch); err != nil {
		return fmt.Errorf("could not rebase branch %s: %w", branch, err)
	}

	return nil
}

func DeleteBranch(branch string) error {
	if _, err := execGit("branch", "-D", branch); err != nil {
		return fmt.Errorf("could not delete branch %s: %w", branch, err)
	}

	return nil
}
