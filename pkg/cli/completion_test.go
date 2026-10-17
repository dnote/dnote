/* Copyright 2026 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var expectedCommands = []string{
	"add",
	"view",
	"edit",
	"remove",
	"find",
	"sync",
	"login",
	"logout",
	"version",
	"help",
}

func getZshCompletionPath(t *testing.T) string {
	t.Helper()
	// Try current directory first (when running from pkg/cli)
	candidates := []string{
		"dnote-completion.zsh",
		filepath.Join("pkg", "cli", "dnote-completion.zsh"),
		filepath.Join("..", "..", "pkg", "cli", "dnote-completion.zsh"),
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			abs, err := filepath.Abs(p)
			if err == nil {
				return abs
			}
			return p
		}
	}

	t.Fatalf("could not find dnote-completion.zsh")
	return ""
}

func getBashCompletionPath(t *testing.T) string {
	t.Helper()
	candidates := []string{
		"dnote-completion.bash",
		filepath.Join("pkg", "cli", "dnote-completion.bash"),
		filepath.Join("..", "..", "pkg", "cli", "dnote-completion.bash"),
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			abs, err := filepath.Abs(p)
			if err == nil {
				return abs
			}
			return p
		}
	}

	t.Fatalf("could not find dnote-completion.bash")
	return ""
}

func TestZshCompletionHeaderAndAutoloadConfig(t *testing.T) {
	contentBytes, err := os.ReadFile(getZshCompletionPath(t))
	if err != nil {
		t.Fatalf("failed to read zsh completion file: %v", err)
	}
	content := string(contentBytes)
	lines := strings.Split(content, "\n")

	if len(lines) == 0 {
		t.Fatalf("zsh completion file is empty")
	}

	// Line 1 must be #compdef dnote
	if strings.TrimSpace(lines[0]) != "#compdef dnote" {
		t.Errorf("expected first line to be '#compdef dnote', got '%s'", lines[0])
	}

	// Must define function _dnote
	if !strings.Contains(content, "_dnote() {") && !strings.Contains(content, "function _dnote") {
		t.Errorf("completion script must define _dnote() function")
	}

	// Must contain compdef registration for sourcing
	if !strings.Contains(content, "compdef _dnote dnote") {
		t.Errorf("completion script must contain 'compdef _dnote dnote' registration")
	}

	// Must contain autoload execution guard using funcstack
	if !strings.Contains(content, "$funcstack[1]") || !strings.Contains(content, "_dnote") {
		t.Errorf("completion script must include autoload guard checking funcstack")
	}

	// Internal helper must be _get_booknames
	if !strings.Contains(content, "_get_booknames()") && !strings.Contains(content, "function _get_booknames") {
		t.Errorf("internal helper function should be named _get_booknames")
	}
}

func TestZshCompletionCommandCoverage(t *testing.T) {
	contentBytes, err := os.ReadFile(getZshCompletionPath(t))
	if err != nil {
		t.Fatalf("failed to read zsh completion file: %v", err)
	}
	content := string(contentBytes)

	for _, cmd := range expectedCommands {
		expectedPrefix := fmt.Sprintf("'%s:", cmd)
		if !strings.Contains(content, expectedPrefix) {
			t.Errorf("zsh completion missing command '%s' in _1st_arguments", cmd)
		}
	}

	// Check book name completion cases
	if !strings.Contains(content, "v|view|a|add") {
		t.Errorf("zsh completion missing book name completion case pattern 'v|view|a|add'")
	}
}

func TestZshCompletionSyntax(t *testing.T) {
	zshPath, err := exec.LookPath("zsh")
	if err != nil {
		t.Skip("zsh not found in PATH, skipping zsh syntax parser test")
	}

	scriptPath := getZshCompletionPath(t)
	cmd := exec.Command(zshPath, "-n", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("zsh syntax verification failed: %v\nOutput: %s", err, string(output))
	}
}

func TestZshCompletionSourcing(t *testing.T) {
	zshPath, err := exec.LookPath("zsh")
	if err != nil {
		t.Skip("zsh not found in PATH, skipping zsh sourcing test")
	}

	scriptPath := getZshCompletionPath(t)
	zshCmd := fmt.Sprintf(`
autoload -Uz compinit
compinit -u -D
source %s

if [[ "${_comps[dnote]}" != "_dnote" ]]; then
  echo "FAIL: expected _comps[dnote] to be '_dnote', got '${_comps[dnote]}'"
  exit 1
fi

if [[ "$(whence -w _dnote)" != "_dnote: function" ]]; then
  echo "FAIL: expected _dnote to be a function, got '$(whence -w _dnote)'"
  exit 1
fi
`, scriptPath)

	cmd := exec.Command(zshPath, "-f", "-c", zshCmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("zsh sourcing test failed: %v\nOutput: %s", err, string(output))
	}
}

func TestZshCompletionAutoloadViaFpath(t *testing.T) {
	zshPath, err := exec.LookPath("zsh")
	if err != nil {
		t.Skip("zsh not found in PATH, skipping zsh fpath autoload test")
	}

	scriptPath := getZshCompletionPath(t)
	tmpDir := t.TempDir()
	fpathTarget := filepath.Join(tmpDir, "_dnote")

	scriptBytes, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("failed to read script: %v", err)
	}
	if err := os.WriteFile(fpathTarget, scriptBytes, 0644); err != nil {
		t.Fatalf("failed to write fpath file: %v", err)
	}

	zshCmd := fmt.Sprintf(`
fpath=(%s $fpath)
autoload -Uz compinit
compinit -u -D

if [[ "${_comps[dnote]}" != "_dnote" ]]; then
  echo "FAIL: expected _comps[dnote] to be '_dnote' via fpath, got '${_comps[dnote]}'"
  exit 1
fi

# Simulate completion trigger
words=(dnote "")
CURRENT=2
"${_comps[dnote]}"

if [[ "$(whence -w _dnote)" != "_dnote: function" ]]; then
  echo "FAIL: expected _dnote to be a function after invocation, got '$(whence -w _dnote)'"
  exit 1
fi
`, tmpDir)

	cmd := exec.Command(zshPath, "-f", "-c", zshCmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("zsh autoload test failed: %v\nOutput: %s", err, string(output))
	}
}

func TestBashCompletionSyntaxAndCoverage(t *testing.T) {
	bashPath, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found in PATH, skipping bash syntax test")
	}

	scriptPath := getBashCompletionPath(t)
	cmd := exec.Command(bashPath, "-n", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("bash syntax verification failed: %v\nOutput: %s", err, string(output))
	}

	contentBytes, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("failed to read bash completion file: %v", err)
	}
	content := string(contentBytes)

	for _, c := range expectedCommands {
		expectedWord := fmt.Sprintf(`"%s"`, c)
		if !strings.Contains(content, expectedWord) {
			t.Errorf("bash completion missing command %s", expectedWord)
		}
	}
}
