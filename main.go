package main

import (
	"fmt"
	"os"

	"ytui/internal/ui"
	"ytui/internal/ytdlp"
)

func main() {
	if err := ytdlp.CheckInstalled(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n\nPlease install yt-dlp first:\n  pacman -S yt-dlp\n  # or\n  pip install yt-dlp\n", err)
		os.Exit(1)
	}

	if err := ui.RunWithClipboard(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
