package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [repository-url]",
	Short: "Initialize a new git-manager workspace",
	Long: `Initialize a new git-manager workspace with a git repository.
This command will clone the repository and set up the initial worktree structure.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := args[0]
		initWorkspace(repoURL)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initWorkspace(repoURL string) {
	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Extract repository name from URL
	repoName := filepath.Base(repoURL)
	if len(repoName) > 4 && repoName[len(repoName)-4:] == ".git" {
		repoName = repoName[:len(repoName)-4]
	}

	// Create directory structure
	repoDir := filepath.Join(currentDir, repoName)
	mainDir := filepath.Join(repoDir, "main")

	// Create directories
	if err := os.MkdirAll(mainDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directories: %v\n", err)
		os.Exit(1)
	}

	// Clone the repository
	fmt.Printf("Cloning repository %s...\n", repoURL)
	cloneCmd := exec.Command("git", "clone", "--bare", repoURL, filepath.Join(repoDir, ".git"))
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error cloning repository: %v\n", err)
		os.Exit(1)
	}

	// Create initial worktree
	fmt.Println("Creating initial worktree...")
	worktreeCmd := exec.Command("git", "-C", filepath.Join(repoDir, ".git"), "worktree", "add", mainDir)
	worktreeCmd.Stdout = os.Stdout
	worktreeCmd.Stderr = os.Stderr
	if err := worktreeCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating worktree: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nGit Manager workspace initialized successfully in %s\n", repoDir)
	fmt.Printf("Main worktree created at %s\n", mainDir)
	fmt.Println("\nYou can now cd into the main directory and start working:")
	fmt.Printf("  cd %s\n", mainDir)
}
