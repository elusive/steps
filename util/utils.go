package util

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

/**
 * function to find files by ext.
 * e.g.
 *  for _, s := range find("/root", ".md") {
 *     println(s)
 *  }
 */
func Find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}

/**
 * Returns true if file exists, otherwise false.
 */
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

/**
 * Returns process id for the provided program name.
 */
func KillProcessAsync(name string) {
	time.Sleep(3 * time.Second)
	go func() {
		// Use WMIC to find and kill notepad.exe processes
		cmd := exec.Command("wmic", "process", "where", fmt.Sprintf("name='%s'", name), "get", "processid")
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error executing WMIC:", err)
		}

		// Parse the output to extract process IDs
		processIDs := extractProcessIDs(string(output))

		// Kill each process
		for _, pid := range processIDs {
			killCmd := exec.Command("taskkill", "/F", "/PID", pid)
			if err := killCmd.Run(); err != nil {
				fmt.Printf("Error killing process %s: %s\n", pid, err)
			} else {
				fmt.Printf("Killed process %s\n", pid)
			}
		}
	}()
}

func extractProcessIDs(output string) []string {
	lines := strings.Split(output, "\n")
	var processIDs []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "ProcessId" && line != "" {
			processIDs = append(processIDs, line)
		}
	}
	return processIDs
}
