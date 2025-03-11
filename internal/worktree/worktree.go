package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Repository struct {
	Name string
	Path string
}

// Info represents information about a git worktree
type Info struct {
	Repository Repository
	Path       string
	Branch     string
	Commit     string
	IsBare     bool
}

// GetWorktreeInfo returns information about all worktrees in the repository
// dir can be a .git directory or anywhere `git` commands can be run
func GetWorktreeInfo(dir string) ([]Info, error) {
	// Run git worktree list command with porcelain output
	cmd := exec.Command("git", "-C", dir, "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error listing worktrees: %v", err)
	}

	// Parse the output
	var worktrees []Info
	var currentWorktree *Info

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if currentWorktree != nil {
				worktrees = append(worktrees, *currentWorktree)
				currentWorktree = nil
			}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			if currentWorktree != nil {
				worktrees = append(worktrees, *currentWorktree)
			}
			path := strings.TrimPrefix(line, "worktree ")
			currentWorktree = &Info{
				Path: path,
			}
		} else if currentWorktree != nil {
			if strings.HasPrefix(line, "branch ") {
				branch := strings.TrimPrefix(line, "branch ")
				// The branch is usually in the format "refs/heads/branch-name"
				parts := strings.Split(branch, "/")
				if len(parts) > 0 {
					currentWorktree.Branch = parts[len(parts)-1]
				} else {
					currentWorktree.Branch = branch
				}
			} else if strings.HasPrefix(line, "HEAD ") {
				currentWorktree.Commit = strings.TrimPrefix(line, "HEAD ")
			} else if strings.HasPrefix(line, "bare") {
				currentWorktree.IsBare = true
			}
		}
	}

	if currentWorktree != nil {
		worktrees = append(worktrees, *currentWorktree)
	}

	return worktrees, nil
}

// IsGitRepository checks if the given directory is a git repository
func IsGitRepository(dir string) bool {
	return IsBareRepository(dir) || IsWorktree(dir)
}

func IsBareRepository(dir string) bool {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--is-bare-repository")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

// IsWorktree checks if the given directory is a git worktree
func IsWorktree(dir string) bool {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

// FindGitDir attempts to find the .git directory by traversing up the directory tree
func FindGitDir(startDir string) (string, error) {
	currentDir := startDir

	// First, check if we're in a worktree
	gitFile := filepath.Join(currentDir, ".git")
	if info, err := os.Stat(gitFile); err == nil {
		// If .git is a file, it's a worktree reference
		if !info.IsDir() {
			content, err := os.ReadFile(gitFile)
			if err != nil {
				return "", fmt.Errorf("error reading .git file: %v", err)
			}

			// Parse the gitdir reference
			gitdirLine := string(content)
			if strings.HasPrefix(gitdirLine, "gitdir: ") {
				gitdirPath := strings.TrimSpace(strings.TrimPrefix(gitdirLine, "gitdir: "))

				// The gitdir path might be relative, so resolve it
				if !filepath.IsAbs(gitdirPath) {
					gitdirPath = filepath.Join(currentDir, gitdirPath)
				}

				// The path usually points to .git/worktrees/branch, we need to go up two levels
				return filepath.Dir(filepath.Dir(gitdirPath)), nil
			}
		} else {
			// If .git is a directory, we're in the main repository
			return gitFile, nil
		}
	}

	// If we're not in a worktree, traverse up the directory tree
	for {
		// Check if we've reached the root directory
		if currentDir == filepath.Dir(currentDir) {
			return "", fmt.Errorf("not in a git repository or git-manager workspace")
		}

		// Move up one directory
		currentDir = filepath.Dir(currentDir)

		// Check if the current directory contains a .git directory
		gitDir := filepath.Join(currentDir, ".git")
		if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
			return gitDir, nil
		}
	}
}
