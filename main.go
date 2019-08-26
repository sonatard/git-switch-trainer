package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		const checkout = "checkout"
		checkoutAliasConfigs, exitCode, err := grepGitConfig(checkout)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(os.Stderr, "Error: failed to read alias from git config")
			os.Exit(exitCode)
		}

		targets := append(aliasFromConfigs(checkoutAliasConfigs), checkout)

		subCommand := os.Args[1]
		if matchSubCommand(subCommand, targets) {
			fmt.Fprintf(os.Stderr, "Error: Use git switch or git restore instead of git checkout\n\n")
			os.Exit(1)
		}
	}

	if exitCode, err := gitExec(os.Args[1:]); err != nil {
		os.Exit(exitCode)
	}
}

func grepGitConfig(target string) ([]string, int, error) {
	var stdout bytes.Buffer
	cmd := exec.Command("git", "config", "--list")
	cmd.Stdout = &stdout
	cmd.Stderr = ioutil.Discard
	if err := cmd.Run(); err != nil {
		exitCode := cmd.ProcessState.ExitCode()
		return nil, exitCode, err
	}

	var matches []string
	configLines := strings.Fields(stdout.String())
	for _, configLine := range configLines {
		if strings.Contains(configLine, target) {
			matches = append(matches, configLine)
		}
	}

	return matches, 0, nil
}

// alias.co=checkout -> co
func aliasFromConfigs(aliasConfigs []string) []string {
	var aliases []string
	for _, aliasConfig := range aliasConfigs {
		// alias.co=checkout -> [alias.co checkout]
		aliasConfigSplit := strings.Split(aliasConfig, "=")
		if len(aliasConfigSplit) > 1 {
			// alias.co -> [alias co]
			aliasConfigSplit2 := strings.Split(aliasConfigSplit[0], ".")
			if len(aliasConfigSplit2) > 1 {
				// co
				aliases = append(aliases, aliasConfigSplit2[1])
			}
		}
	}

	return aliases
}

func matchSubCommand(subCommand string, targets []string) bool {
	for _, target := range targets {
		if target == subCommand {
			return true
		}
	}

	return false
}

func gitExec(args []string) (int, error) {
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		exitCode := cmd.ProcessState.ExitCode()
		return exitCode, err
	}

	return 0, nil
}
