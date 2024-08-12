package main

import (
	"bufio"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

var (
	fallbackHostsPath   = getHostsPath()
	fallbackSitesFile   = "sites.txt"
	fallbackWebsiteList []string
)

func fallbackMain() {
	loadFallbackWebsites()
	unblockFallbackWebsites()
}

func getHostsPath() string {
	switch runtime.GOOS {
	case "windows":
		return "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "linux", "darwin":
		return "/etc/hosts"
	default:
		color.Red("Unsupported operating system: %s", runtime.GOOS)
		os.Exit(1)
		return ""
	}
}

func loadFallbackWebsites() {
	file, err := os.Open(fallbackSitesFile)
	if err != nil {
		color.Red("Error opening %s: %v", fallbackSitesFile, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		website := strings.TrimSpace(scanner.Text())
		if website != "" {
			fallbackWebsiteList = append(fallbackWebsiteList, website)
		}
	}

	if err := scanner.Err(); err != nil {
		color.Red("Error reading %s: %v", fallbackSitesFile, err)
		os.Exit(1)
	}
}

func unblockFallbackWebsites() {
	file, err := os.Open(fallbackHostsPath)
	if err != nil {
		color.Red("Error opening hosts file: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		keepLine := true
		for _, website := range fallbackWebsiteList {
			if strings.Contains(line, website) {
				keepLine = false
				break
			}
		}
		if keepLine {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		color.Red("Error scanning hosts file: %v", err)
		os.Exit(1)
	}

	err = os.WriteFile(fallbackHostsPath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		color.Red("Error writing to hosts file: %v", err)
	} else {
		color.Green("Successfully unblocked all specified websites.")
	}
}
