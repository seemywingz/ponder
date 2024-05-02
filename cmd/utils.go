package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/pterm/pterm"
)

var spinner *pterm.SpinnerPrinter
var moonSequence = []string{"🌑 ", "🌒 ", "🌓 ", "🌔 ", "🌕 ", "🌖 ", "🌗 ", "🌘 "}
var ponderSpinner = &pterm.SpinnerPrinter{
	Sequence:            []string{"▀ ", " ▀", " ▄", "▄ "},
	Style:               &pterm.ThemeDefault.SpinnerStyle,
	Delay:               time.Millisecond * 200,
	ShowTimer:           false,
	TimerRoundingFactor: time.Second,
	TimerStyle:          &pterm.ThemeDefault.TimerStyle,
	MessageStyle:        &pterm.ThemeDefault.SpinnerTextStyle,
	InfoPrinter:         &pterm.Info,
	SuccessPrinter:      &pterm.Success,
	FailPrinter:         &pterm.Error,
	WarningPrinter:      &pterm.Warning,
	RemoveWhenDone:      true,
	Text:                "Pondering...",
}

func expanding(emoji string, maxRadius int) []string {
	totalLength := maxRadius*2 - 1          // Total fixed length for each line
	sequence := make([]string, maxRadius*2) // Frames for expanding and contracting

	// Generate expanding sequence
	for i := 0; i < maxRadius; i++ {
		spacesBefore := strings.Repeat(" ", maxRadius-i-1)
		numEmojis := i + 1
		emojis := strings.Repeat(emoji+" ", numEmojis)
		spacesAfterCount := totalLength - len(spacesBefore) - len(emojis)
		if spacesAfterCount < 0 {
			spacesAfterCount = 0 // Ensure spacesAfterCount is never negative
		}
		spacesAfter := strings.Repeat(" ", spacesAfterCount)
		sequence[i] = spacesBefore + emojis + spacesAfter
	}

	// Generate contracting sequence
	for i := 0; i < maxRadius; i++ {
		spacesBefore := strings.Repeat(" ", i+1)
		numEmojis := maxRadius - i
		emojis := strings.Repeat(emoji+" ", numEmojis)
		spacesAfterCount := totalLength - len(spacesBefore) - len(emojis)
		if spacesAfterCount < 0 {
			spacesAfterCount = 0 // Prevent negative space count
		}
		spacesAfter := strings.Repeat(" ", spacesAfterCount)
		sequence[maxRadius+i] = spacesBefore + emojis + spacesAfter
	}

	return sequence
}

func prettyPrint(message string) {
	lines := strings.Split(message, "\n")
	var codeBuffer bytes.Buffer
	var inCodeBlock bool
	var currentLexer chroma.Lexer

	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}
	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	// Regex to find inline code and double-quoted text
	backtickRegex := regexp.MustCompile("`([^`]*)`")
	doubleQuoteRegex := regexp.MustCompile(`"([^"]*)"`)
	cyan := "\033[36m"   // Cyan color ANSI escape code
	yellow := "\033[33m" // Yellow color ANSI escape code
	reset := "\033[0m"   // Reset ANSI escape code

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			if inCodeBlock {
				// Ending a code block, apply syntax highlighting
				iterator, err := currentLexer.Tokenise(nil, codeBuffer.String())
				if err == nil {
					formatter.Format(os.Stdout, style, iterator)
				}
				fmt.Println() // Ensure there's a newline after the code block
				codeBuffer.Reset()
				inCodeBlock = false
			} else {
				// Starting a code block
				inCodeBlock = true
				lang := strings.TrimPrefix(strings.TrimSpace(line), "```")
				currentLexer = lexers.Get(lang)
				if currentLexer == nil {
					currentLexer = lexers.Fallback
				}
				continue // Skip the line with opening backticks
			}
		} else if inCodeBlock {
			codeBuffer.WriteString(line + "\n") // Collect code lines
		} else {
			// Process and set colors
			processedLine := line
			processedLine = backtickRegex.ReplaceAllStringFunc(processedLine, func(match string) string {
				return cyan + strings.Trim(match, "`") + reset
			})
			processedLine = doubleQuoteRegex.ReplaceAllStringFunc(processedLine, func(match string) string {
				return yellow + match + reset
			})
			fmt.Println("    " + processedLine) // Print with white color
		}
	}

	// Flush the remaining content if still in a code block
	if inCodeBlock {
		iterator, err := currentLexer.Tokenise(nil, codeBuffer.String())
		if err == nil {
			formatter.Format(os.Stdout, style, iterator)
		}
		fmt.Println() // Ensure there's a newline after the code block
	}
}

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

// playAudio plays audio from a byte slice.
func playAudio(audioContent []byte) {
	if verbose {
		fmt.Println("🔊 Playing audio...")
	}

	// Create an io.Reader from the byte slice
	reader := bytes.NewReader(audioContent)

	// Wrap the reader in a NopCloser to make it an io.ReadCloser.
	readCloser := io.NopCloser(reader)

	// Decode the MP3 stream.
	streamer, format, err := mp3.Decode(readCloser)
	catchErr(err)
	defer streamer.Close()

	// Initialize the speaker with the sample rate of the audio and a buffer size.
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	catchErr(err)

	// Play the decoded audio.
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	// Wait for the audio to finish playing.
	<-done
}

func playMP3File(file string) {
	if verbose {
		fmt.Println("🔊 Playing audio file:", file)
	}

	f, err := os.Open(file)
	catchErr(err)
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	catchErr(err)
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	catchErr(err)

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}
