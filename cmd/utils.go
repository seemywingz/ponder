package cmd

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"runtime"
)

func catchErr(err error) {
	if err != nil {
		fmt.Println("💔", err)
		// os.Exit(1)
	}
}

func formatPrompt(prompt string) string {
	// Replace any characters that are not letters, numbers, or underscores with dashes
	return regexp.MustCompile(`[^a-zA-Z0-9_]+`).ReplaceAllString(prompt, "-")
}

func fileNameFromURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	catchErr(err)
	// Get the last path component of the URL
	filename := filepath.Base(u.Path)
	// Replace any characters that are not letters, numbers, or underscores with dashes
	filename = regexp.MustCompile(`[^a-zA-Z0-9_]+`).ReplaceAllString(filename, "-")
	// Limit the filename to 255 characters
	if len(filename) >= 255 {
		filename = filename[:255]
	}
	return filename
}

func trace() {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("%s:%d\n%s\n", file, line, f.Name())
}
