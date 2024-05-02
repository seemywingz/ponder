package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/pterm/pterm"
)

var spinner *pterm.SpinnerPrinter
var moonSpinner = &pterm.SpinnerPrinter{
	Sequence:            []string{"🌑 ", "🌒 ", "🌓 ", "🌔 ", "🌕 ", "🌖 ", "🌗 ", "🌘 "},
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
}

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
