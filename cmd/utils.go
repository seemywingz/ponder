package cmd

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
)

func catchErr(err error, level ...string) {
	if err != nil {
		// Default level is "warn" if none is provided
		lvl := "warn"
		if len(level) > 0 {
			lvl = level[0] // Use the provided level
		}

		switch lvl {
		case "warn":
			fmt.Println("💔 Warning:", err)
		case "fatal":
			fmt.Println("💀 Fatal:", err)
			os.Exit(1)
		}
	}
}

func formatPrompt(prompt string) string {
	// Replace any characters that are not letters, numbers, or underscores with dashes
	return regexp.MustCompile(`[^a-zA-Z0-9_]+`).ReplaceAllString(prompt, "-")
}

func trace() {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("%s:%d\n%s\n", file, line, f.Name())
}
