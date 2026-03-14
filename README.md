# ytui

A terminal UI for downloading videos and audio using yt-dlp.

## Features

- Clipboard integration - automatically detects URLs in clipboard
- Interactive TUI built with Bubble Tea and huh
- Video/Audio download selection
- Resolution selection: 720p, 1080p, 1440p, 4K
- Format selection: MP4, MKV, WebM (video) / MP3, M4A, FLAC, WAV (audio)
- Custom download folder
- Desktop notifications on completion/error
- Automatic fallback if requested resolution unavailable
- Embeds thumbnail and metadata

## Requirements

- [yt-dlp](https://github.com/yt-dlp/yt-dlp) installed
- [FFmpeg](https://ffmpeg.org/) for merging video+audio and audio conversion
- `notify-send` for desktop notifications (usually comes with libnotify)
- Clipboard tools: `wl-paste`, `xclip`, or `xsel` (for clipboard detection)

### Arch Linux Installation

```bash
pacman -S yt-dlp ffmpeg libnotify
# For clipboard support (choose one):
pacman -S wl-clipboard      # Wayland
pacman -S xclip             # X11
pacman -S xsel              # X11 alternative
```

## Installation

### Build from Source

```bash
# Clone and build
git clone https://github.com/yourusername/ytui.git
cd ytui
go build -o ytui

# Or install globally
go install
```

### Build Flags

For a smaller binary:
```bash
go build -ldflags="-s -w" -o ytui
```

For a static binary (requires musl-gcc):
```bash
CGO_ENABLED=1 CC=musl-gcc go build -ldflags="-linkmode external -extldflags -static" -o ytui
```

## Usage

1. Copy a video URL to your clipboard
2. Run `ytui`
3. If URL is detected in clipboard, it will be pre-filled
4. Follow the prompts:
   - Confirm URL (or enter manually if not detected)
   - Choose Video or Audio
   - Select resolution (for video)
   - Select format
   - Choose download folder (default: ~/Downloads)
   - Press Enter to confirm and download

### Keyboard Shortcuts

- `Enter` - Confirm selection / Start download
- `↑/↓` or `j/k` - Navigate options
- `Ctrl+C` - Cancel/Quit
- `R` - Retry after error (in error state)

## Default Values

| Prompt | Default |
|--------|---------|
| Type | Video |
| Resolution | 1080p |
| Video Format | MP4 |
| Audio Format | MP3 |
| Download Folder | ~/Downloads |

## Configuration

yt-dlp config file (optional): `~/.config/yt-dlp/config`

Example:
```
--format bv*+ba
--merge-output-format mp4
--embed-thumbnail
--add-metadata
```

## Troubleshooting

### "yt-dlp is not installed"
Make sure yt-dlp is in your PATH:
```bash
which yt-dlp
```

### No clipboard detection
Install one of: `wl-clipboard`, `xclip`, or `xsel`

### No notifications
Make sure `notify-send` works:
```bash
notify-send "Test" "Hello"
```

