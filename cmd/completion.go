package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/engmtcdrm/mellon/app"
)

// detectShell attempts to detect the current user's shell.
func detectShell() string {
	shell := os.Getenv("SHELL")

	if strings.HasSuffix(shell, "zsh") {
		return "zsh"
	} else if strings.HasSuffix(shell, "bash") {
		return "bash"
	} else if strings.HasSuffix(shell, "fish") {
		return "fish"
	} else if isPowerShell() {
		return "powershell"
	}

	return ""
}

// isPowerShell checks if the current shell is PowerShell.
func isPowerShell() bool {
	// Check for PowerShell-specific environment variables
	if os.Getenv("PSVersionTable") != "" || os.Getenv("PSModulePath") != "" || os.Getenv("PROFILE") != "" {
		return true
	}

	return false
}

// completionFilePath returns the expected path for the shell completion file
// based on the shell type and user's home directory.
func completionFilePath(shell, homeDir string) string {
	switch shell {
	case "bash":
		return filepath.Join(homeDir, ".local", "share", "bash-completion", "completions", app.Name)
	case "zsh":
		return filepath.Join(homeDir, ".zsh", "completions", fmt.Sprintf("_%s", app.Name))
	case "fish":
		return filepath.Join(homeDir, ".config", "fish", "completions", fmt.Sprintf("%s.fish", app.Name))
	case "powershell":
		profile, _ := exec.Command("powershell", "-Command", "$PROFILE").Output()
		if string(profile) != "" {
			profileDir := filepath.Dir(string(profile))
			return filepath.Join(profileDir, fmt.Sprintf("%s.ps1", app.Name))
		}

		return ""
	default:
		return ""
	}
}

// findInFile checks if a specific line exists in a file.
func findInFile(filePath, searchTerm string) (bool, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == strings.TrimSpace(searchTerm) {
			return true, nil
		}
	}

	return false, scanner.Err()
}

// genZshCompletion generates zsh completion and appends necessary configurations to .zshrc
func genZshCompletion(file *os.File, homeDir string) {
	if err := rootCmd.GenZshCompletion(file); err != nil {
		fmt.Printf("Failed to generate zsh completion script: %v\n", err)
		return
	}

	zshrcFile, err := os.OpenFile(filepath.Join(homeDir, ".zshrc"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open .zshrc file: %v\n", err)
		return
	}
	defer zshrcFile.Close()

	fpath := "fpath=(~/.zsh/completions $fpath)\n"
	foundFpath, err := findInFile(filepath.Join(homeDir, ".zshrc"), fpath)
	if err != nil {
		fmt.Printf("Failed to read .zshrc file: %v\n", err)
		return
	}

	if !foundFpath {
		if _, err := zshrcFile.WriteString("# mellon shell completion\n" + fpath); err != nil {
			fmt.Printf("Failed to write to .zshrc file: %v\n", err)
			return
		}
	}

	autoload := "autoload -U compinit && compinit\n"
	foundAutoload, err := findInFile(filepath.Join(homeDir, ".zshrc"), autoload)
	if err != nil {
		fmt.Printf("Failed to read .zshrc file: %v\n", err)
		return
	}

	if !foundAutoload {
		if _, err := zshrcFile.WriteString(autoload); err != nil {
			fmt.Printf("Failed to write to .zshrc file: %v\n", err)
			return
		}
	}
}

// genPowerShellCompletion generates PowerShell completion and appends necessary configurations to the PowerShell profile
func genPowerShellCompletion(file *os.File) {
	profileFilePathByte, _ := exec.Command("powershell", "-Command", "$PROFILE").Output()
	profileFilePath := strings.TrimSpace(string(profileFilePathByte))

	if profileFilePath == "" {
		fmt.Printf("Could not determine PowerShell profile path.\n")
		return
	}

	if err := rootCmd.GenPowerShellCompletion(file); err != nil {
		fmt.Printf("Failed to generate powershell completion script: %v\n", err)
		return
	}

	profileDir := filepath.Dir(profileFilePath)
	compSource := ". " + filepath.Join(profileDir, fmt.Sprintf("%s.ps1", app.Name))

	profileFile, err := os.OpenFile(profileFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open PowerShell profile file: %v\n", err)
		return
	}
	defer profileFile.Close()

	foundCompPath, err := findInFile(profileFilePath, compSource)
	if err != nil {
		fmt.Printf("Failed to read PowerShell profile file: %v\n", err)
		return
	}

	if !foundCompPath {
		if _, err := profileFile.WriteString("# mellon shell completion\n" + compSource); err != nil {
			fmt.Printf("Failed to write to PowerShell profile file: %v\n", err)
			return
		}
	}
}

// initShellCompletion initializes shell completion for the detected shell.
func initShellCompletion(homeDir string) {
	shell := detectShell()
	completionPath := completionFilePath(shell, homeDir)

	if shell == "" || completionPath == "" {
		return
	}

	// Check if the completion file already exists
	if _, err := os.Stat(completionPath); os.IsNotExist(err) {
		// Create parent directories if they don't exist
		completionDir := filepath.Dir(completionPath)
		if _, err := os.Stat(completionDir); os.IsNotExist(err) {
			err = os.MkdirAll(completionDir, dirMode)
			if err != nil {
				fmt.Printf("Failed to create directory for shell completion: %v\n", err)
				return
			}
		}

		// Create the completion file
		file, err := os.Create(completionPath)
		if err != nil {
			fmt.Printf("Failed to create shell completion file: %v\n", err)
			return
		}
		defer file.Close()

		switch shell {
		case "bash":
			if err := rootCmd.GenBashCompletion(file); err != nil {
				fmt.Printf("Failed to generate bash completion script: %v\n", err)
				return
			}
		case "zsh":
			genZshCompletion(file, homeDir)
		case "fish":
			if err := rootCmd.GenFishCompletion(file, true); err != nil {
				fmt.Printf("Failed to generate fish completion script: %v\n", err)
				return
			}
		case "powershell":
			genPowerShellCompletion(file)
		}

		// fmt.Printf("Shell completion script created at %s. Please source it in your shell configuration file.\n", pp.Green(completionPath))
	}
}
