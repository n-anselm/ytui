package ytdlp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Options struct {
	URL             string
	IsAudio         bool
	Resolution      string
	Format          string
	DownloadFolder  string
}

func CheckInstalled() error {
	cmd := exec.Command("yt-dlp", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlp is not installed or not in PATH")
	}
	return nil
}

func (o *Options) BuildArgs() ([]string, error) {
	var args []string

	if o.IsAudio {
		args = []string{
			o.URL,
			"-x",
			"--audio-format", o.Format,
		}
	} else {
		heightFilter := o.resolutionToHeightFilter()
		mergeFormat := o.Format
		if mergeFormat == "webm" {
			mergeFormat = "webm"
		} else {
			mergeFormat = "mp4"
		}

		args = []string{
			o.URL,
			"-f", fmt.Sprintf("bestvideo[height<=%s]+bestaudio/best", heightFilter),
			"--merge-output-format", mergeFormat,
		}
	}

	args = append(args,
		"--embed-thumbnail",
		"--add-metadata",
		"-o", filepath.Join(o.DownloadFolder, "%(title)s.%(ext)s"),
	)

	return args, nil
}

func (o *Options) resolutionToHeightFilter() string {
	switch o.Resolution {
	case "720p":
		return "720"
	case "1080p":
		return "1080"
	case "1440p":
		return "1440"
	case "4K":
		return "2160"
	default:
		return "1080"
	}
}

func (o *Options) Run() (string, string, error) {
	args, err := o.BuildArgs()
	if err != nil {
		return "", "", err
	}

	cmd := exec.Command("yt-dlp", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", "", err
	}

	return o.URL, o.DownloadFolder, nil
}

func GetVideoTitle(url string) (string, error) {
	cmd := exec.Command("yt-dlp", "--print", "title", "--no-download", url)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func IsPlaylist(url string) (bool, error) {
	cmd := exec.Command("yt-dlp", "--dump-json", "--no-download", url)
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	jsonStr := string(out)
	if strings.Contains(jsonStr, `"playlist_count"`) || strings.Contains(jsonStr, `"entries"`) {
		return true, nil
	}
	return false, nil
}
