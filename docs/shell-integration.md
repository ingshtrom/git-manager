# Shell Integration Guide

This document explains how Git Manager's shell integration works and how to extend it to other commands.

## Overview

Git Manager uses a special output format to communicate with the shell wrapper. This allows commands to perform actions that would normally be impossible for a child process, such as changing the parent shell's current directory.

## How It Works

1. Git Manager commands output lines with a special prefix: `git-manager-eval:`
2. The shell wrapper captures this output and evaluates any commands with this prefix
3. All other output is displayed normally

This approach allows Git Manager to:
- Change directories (`cd`)
- Set environment variables
- Execute other shell commands
- Perform any action that would normally require shell integration

## Supported Shell Commands

Here are some examples of shell commands that can be executed using the `git-manager-eval:` prefix:

```
git-manager-eval:cd "/path/to/directory"
git-manager-eval:export VARIABLE="value"
git-manager-eval:source .env
```

## Extending to Other Commands

To add shell integration to a new command, follow these steps:

1. Identify where shell interaction is needed
2. Output the shell command with the `git-manager-eval:` prefix
3. Provide fallback instructions for users without shell integration

### Example: Adding Shell Integration to a Command

```go
func myCommand() {
    // Command logic...
    
    // Output a shell command to be evaluated
    fmt.Printf("git-manager-eval:cd %q\n", somePath)
    
    // Provide fallback instructions
    fmt.Println("\nIf you're not using shell integration, run:")
    fmt.Printf("  cd %s\n", somePath)
}
```

## Best Practices

When using shell integration:

1. **Always quote paths** to handle spaces and special characters
2. **Provide fallback instructions** for users without shell integration
3. **Keep shell commands simple** to ensure compatibility across different shells
4. **Document the behavior** so users understand what's happening

## Security Considerations

The shell integration executes commands directly in the user's shell. To ensure security:

1. **Never execute user-provided input** without proper validation
2. **Limit commands to necessary operations** like directory changes
3. **Quote all arguments** to prevent shell injection

## Testing Shell Integration

To test shell integration:

1. Install the shell integration for your shell
2. Run a command that uses shell integration
3. Verify that the expected action (e.g., directory change) occurs
4. Test with and without shell integration to ensure fallback instructions work

## Troubleshooting

If shell integration isn't working:

1. Ensure the shell integration is properly installed
2. Check that the command is outputting the `git-manager-eval:` prefix
3. Verify that the shell wrapper is correctly parsing the output
4. Try running the command with verbose output to see what's happening

## Example: Custom Command with Shell Integration

Here's an example of a custom command that uses shell integration to set up a development environment:

```go
func setupDevEnv(projectName string) {
    // Create project directory
    projectPath := filepath.Join(os.Getenv("HOME"), "projects", projectName)
    os.MkdirAll(projectPath, 0755)
    
    // Initialize git repository
    exec.Command("git", "init", projectPath).Run()
    
    // Change to project directory using shell integration
    fmt.Printf("git-manager-eval:cd %q\n", projectPath)
    
    // Set environment variables using shell integration
    fmt.Printf("git-manager-eval:export PROJECT_ROOT=%q\n", projectPath)
    
    // Provide fallback instructions
    fmt.Println("\nIf you're not using shell integration, run:")
    fmt.Printf("  cd %s\n", projectPath)
    fmt.Printf("  export PROJECT_ROOT=%s\n", projectPath)
}
``` 
