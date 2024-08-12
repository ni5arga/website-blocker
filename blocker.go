package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	hostsPath   string
	redirectIP  = "127.0.0.1"
	startBlock  = "09:00"
	endBlock    = "17:00"
	sitesFile   = "sites.txt"
	websiteList []string
)

func init() {
	// Detect the operating system
	switch runtime.GOOS {
	case "windows":
		hostsPath = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "linux", "darwin": // "darwin" is macOS
		hostsPath = "/etc/hosts"
	default:
		color.Red("Unsupported operating system: %s", runtime.GOOS)
		os.Exit(1)
	}
}

func main() {
	loadWebsites()

	for {
		now := time.Now()
		currentTime := now.Format("15:04")
		blockTime, _ := time.Parse("15:04", startBlock)
		unblockTime, _ := time.Parse("15:04", endBlock)
		currentParsedTime, _ := time.Parse("15:04", currentTime)

		if currentParsedTime.After(blockTime) && currentParsedTime.Before(unblockTime) {
			blockWebsites()
		} else {
			unblockWebsites()
		}

		// Check every minute
		time.Sleep(1 * time.Minute)
	}
}

func loadWebsites() {
	file, err := os.Open(sitesFile)
	if err != nil {
		color.Red("Error reading %s: %v", sitesFile, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		website := strings.TrimSpace(scanner.Text())
		if website != "" {
			websiteList = append(websiteList, website)
		}
	}

	if err := scanner.Err(); err != nil {
		color.Red("Error scanning %s: %v", sitesFile, err)
	}
}

func blockWebsites() {
	content, err := os.ReadFile(hostsPath)
	if err != nil {
		color.Red("Error reading hosts file: %v", err)
		return
	}

	for _, website := range websiteList {
		if !strings.Contains(string(content), website) {
			entry := fmt.Sprintf("%s %s\n", redirectIP, website)
			err := appendToFile(hostsPath, entry)
			if err != nil {
				color.Red("Error writing to hosts file: %v", err)
			} else {
				color.Green("Blocked %s", website)
			}
		} else {
			color.Yellow("%s is already blocked.", website)
		}
	}
}

func unblockWebsites() {
	file, err := os.Open(hostsPath)
	if err != nil {
		color.Red("Error opening hosts file: %v", err)
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		blocked := false
		for _, website := range websiteList {
			if strings.Contains(line, website) {
				blocked = true
				break
			}
		}
		if !blocked {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		color.Red("Error scanning hosts file: %v", err)
		return
	}

	err = os.WriteFile(hostsPath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		color.Red("Error writing to hosts file: %v", err)
	} else {
		color.Green("Unblocked all specified websites.")
	}
}

func appendToFile(filename, text string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(text); err != nil {
		return err
	}

	return nil
}
