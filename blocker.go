package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	hostsPath   string
	redirectIP  = "127.0.0.1"
	startBlock  = "09:00"
	endBlock    = "17:00"
	sitesFile   = "sites.txt"
	websiteList []string
	isFallback  bool
)

func init() {
	// OS detection
	switch runtime.GOOS {
	case "windows":
		hostsPath = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "linux", "darwin": // "darwin" is macOS
		hostsPath = "/etc/hosts"
	default:
		fmt.Printf("Unsupported operating system: %s\n", runtime.GOOS)
		os.Exit(1)
	}


	// removing fallback cuz i end up using it as a workaround 

	// fallback
	// flag.BoolVar(&isFallback, "fallback", false, "Run in fallback mode to unblock all sites")
	flag.Parse() 
}

func main() {
	if isFallback {
		// unblockWebsites()
		return
	}

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

		time.Sleep(1 * time.Minute)
	}
}

func loadWebsites() {
	file, err := os.Open(sitesFile)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", sitesFile, err)
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
		fmt.Printf("Error scanning %s: %v\n", sitesFile, err)
	}
}

func blockWebsites() {
	content, err := os.ReadFile(hostsPath)
	if err != nil {
		fmt.Printf("Error reading hosts file: %v\n", err)
		return
	}

	// Create a map to keep track of already blocked websites
	blockedMap := make(map[string]bool)
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, redirectIP) {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				blockedMap[parts[1]] = true
			}
		}
	}

	for _, website := range websiteList {
		if _, exists := blockedMap[website]; !exists {
			entry := fmt.Sprintf("%s %s\n", redirectIP, website)
			err := appendToFile(hostsPath, entry)
			if err != nil {
				fmt.Printf("Error writing to hosts file: %v\n", err)
			} else {
				fmt.Printf("Blocked %s\n", website)
			}
		} else {
			fmt.Printf("%s is already blocked.\n", website)
		}
	}
}

func unblockWebsites() {
	file, err := os.Open(hostsPath)
	if err != nil {
		fmt.Printf("Error opening hosts file: %v\n", err)
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
		fmt.Printf("Error scanning hosts file: %v\n", err)
		return
	}

	err = os.WriteFile(hostsPath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		fmt.Printf("Error writing to hosts file: %v\n", err)
	} else {
		fmt.Println("Unblocked all specified websites.")
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
