package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell [shell-type]",
	Short: "Generate shell integration scripts",
	Long: `Generate shell integration scripts for different shells.
This command outputs shell functions that can be added to your shell configuration
to enable directory switching and other advanced features.

Supported shell types: sh, bash, zsh, fish, nushell`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shellType := args[0]
		generateShellIntegration(shellType)
	},
}

func init() {
	toolCmd.AddCommand(shellCmd)
}

func generateShellIntegration(shellType string) {
	switch shellType {
	case "sh", "bash", "zsh":
		fmt.Println(`# Git Manager Shell Integration
# Add this to your .bashrc, .zshrc, or .profile file

# Main wrapper function for git-manager
function git-manager() {
  # Capture the output of the command
  local output
  output=$(command git-manager "$@" 2>&1)
  local exit_code=$?
  
  # Process the output line by line
  local eval_cmd=""
  while IFS= read -r line; do
    if [[ "$line" == git-manager-eval:* ]]; then
      # Extract the command to evaluate
      eval_cmd="${line#git-manager-eval:}"
    else
      # Print regular output
      echo "$line"
    fi
  done <<< "$output"
  
  # If we found an eval command, execute it
  if [ -n "$eval_cmd" ]; then
    eval "$eval_cmd"
  fi
  
  return $exit_code
}

# Alias for shorter command
alias gm=git-manager

# Tab completion for git-manager
if command -v git-manager &> /dev/null; then
  _git_manager_completion() {
    local cur prev
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    if [ "$prev" = "switch" ] || [ "$prev" = "remove" ]; then
      # Get list of worktrees for completion
      local worktrees=$(command git-manager list | grep -v "Available worktrees" | awk '{print $1}')
      COMPREPLY=( $(compgen -W "${worktrees}" -- ${cur}) )
      return 0
    fi
    
    # Complete the main commands
    if [ "$COMP_CWORD" -eq 1 ]; then
      COMPREPLY=( $(compgen -W "init create list switch remove shell" -- ${cur}) )
      return 0
    fi
  }
  
  complete -F _git_manager_completion git-manager
  complete -F _git_manager_completion gm
fi`)

	case "fish":
		fmt.Println(`# Git Manager Shell Integration for Fish
# Save this to ~/.config/fish/functions/git-manager.fish

function git-manager
  # Capture the output of the command
  set -l output (command git-manager $argv 2>&1)
  set -l exit_code $status
  
  # Process the output line by line
  set -l eval_cmd ""
  for line in $output
    if string match -q "git-manager-eval:*" -- $line
      # Extract the command to evaluate
      set eval_cmd (string replace "git-manager-eval:" "" -- $line)
    else
      # Print regular output
      echo $line
    end
  end
  
  # If we found an eval command, execute it
  if test -n "$eval_cmd"
    eval $eval_cmd
  end
  
  return $exit_code
end

# Alias for shorter command
alias gm=git-manager

# Add completion for git-manager commands
complete -c git-manager -f -n "__fish_use_subcommand" -a "init create list switch remove shell" -d "Git Manager command"
complete -c git-manager -f -n "__fish_seen_subcommand_from switch remove" -a "(command git-manager list | grep -v 'Available worktrees' | awk '{print \$1}')" -d "Worktree"
complete -c gm -f -n "__fish_use_subcommand" -a "init create list switch remove shell" -d "Git Manager command"
complete -c gm -f -n "__fish_seen_subcommand_from switch remove" -a "(command git-manager list | grep -v 'Available worktrees' | awk '{print \$1}')" -d "Worktree"`)

	case "nushell":
		fmt.Println(`# Git Manager Shell Integration for Nushell
# Save this to your Nushell config file

def git-manager [...args] {
  # Capture the output of the command
  let output = (do { ^git-manager $args } | complete)
  
  # Process the output
  let lines = $output.stdout
  let eval_cmd = ""
  
  # Look for eval commands in the output
  for line in $lines {
    if ($line | str starts-with "git-manager-eval:") {
      $eval_cmd = ($line | str replace "git-manager-eval:" "")
    } else {
      # Print regular output
      echo $line
    }
  }
  
  # If we found an eval command, execute it
  if ($eval_cmd != "") {
    nu -c $eval_cmd
  }
  
  # Return the original exit code
  exit $output.exit_code
}

# Alias for shorter command
alias gm = git-manager

# Add completions for git-manager
def "nu-complete git-manager-commands" [] {
  ["init", "create", "list", "switch", "remove", "shell"]
}

def "nu-complete git-manager-worktrees" [] {
  ^git-manager list | lines | skip 1 | each { |line| $line | split row " " | get 0 }
}

# Register the completions
export extern "git-manager" [
  command: string@"nu-complete git-manager-commands"  # The command to run
  --help(-h)                                         # Show help
]

export extern "git-manager switch" [
  worktree: string@"nu-complete git-manager-worktrees"  # The worktree to switch to
  --help(-h)                                           # Show help
]

export extern "git-manager remove" [
  worktree: string@"nu-complete git-manager-worktrees"  # The worktree to remove
  --help(-h)                                           # Show help
]`)

	default:
		fmt.Fprintf(os.Stderr, "Unsupported shell type: %s\n", shellType)
		fmt.Println("Supported shell types: sh, bash, zsh, fish, nushell")
		os.Exit(1)
	}

	fmt.Println("\n# To install, run:")
	switch shellType {
	case "sh":
		fmt.Println("# echo 'source <(git-manager shell sh)' >> ~/.profile")
	case "bash":
		fmt.Println("# echo 'source <(git-manager shell bash)' >> ~/.bashrc")
	case "zsh":
		fmt.Println("# echo 'source <(git-manager shell zsh)' >> ~/.zshrc")
	case "fish":
		fmt.Println("# git-manager shell fish > ~/.config/fish/functions/git-manager.fish")
	case "nushell":
		fmt.Println("# git-manager shell nushell >> ~/.config/nushell/config.nu")
	}
}
