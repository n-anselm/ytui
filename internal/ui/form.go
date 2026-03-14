package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"
	"ytui/internal/clipboard"
	"ytui/internal/validator"
	"ytui/internal/ytdlp"
)

type model struct {
	form           *huh.Form
	url            string
	isVideo        bool
	resolution     string
	format         string
	folder         string
	downloading    bool
	done           bool
	err            error
	successMsg     string
	skipToDl       bool
	spinnerFrame   int
}

func New(cliURL string) *model {
	m := &model{
		url:        cliURL,
		isVideo:    true,
		resolution: "1080p",
		format:     "mp4",
		folder:     filepath.Join(os.Getenv("HOME"), "Downloads"),
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Video URL").
				Placeholder("https://youtube.com/watch?v=...").
				Value(&m.url).
				Key("url").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("URL is required")
					}
					if !validator.IsValidURL(s) {
						return fmt.Errorf("invalid URL format")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Download video or audio?").
				Affirmative("Video").
				Negative("Audio").
				Value(&m.isVideo).
				Key("isVideo"),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select resolution").
				Options(
					huh.NewOption("720p", "720p"),
					huh.NewOption("1080p (recommended)", "1080p"),
					huh.NewOption("1440p", "1440p"),
					huh.NewOption("4K", "4K"),
				).
				Value(&m.resolution).
				Key("resolution"),
		).WithHideFunc(func() bool { return !m.isVideo }),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select format").
				OptionsFunc(func() []huh.Option[string] {
					if m.isVideo {
						return []huh.Option[string]{
							huh.NewOption("MP4 (recommended)", "mp4"),
							huh.NewOption("MKV", "mkv"),
							huh.NewOption("WebM", "webm"),
						}
					}
					return []huh.Option[string]{
						huh.NewOption("MP3 (recommended)", "mp3"),
						huh.NewOption("M4A", "m4a"),
						huh.NewOption("FLAC", "flac"),
						huh.NewOption("WAV", "wav"),
					}
				}, &m.isVideo).
				Value(&m.format).
				Key("format"),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Download folder").
				Placeholder("~/Downloads").
				Value(&m.folder).
				Key("folder").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("folder is required")
					}
					expanded := os.ExpandEnv(s)
					if _, err := os.Stat(expanded); os.IsNotExist(err) {
						return fmt.Errorf("folder does not exist")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeFunc(huh.ThemeDracula))

	return m
}

func (m *model) Init() tea.Cmd {
	return m.form.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.done && m.err == nil {
			if msg.String() == "r" {
				// Reset to initial state for new download
				// Try to get URL from clipboard again
				if clipContent, err := clipboard.Read(); err == nil && validator.IsValidURL(clipContent) {
					m.url = clipContent
				} else {
					m.url = ""
				}
				m.done = false
				m.successMsg = ""
				m.err = nil
				m.form = huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Video URL").
							Placeholder("https://youtube.com/watch?v=...").
							Value(&m.url).
							Key("url").
							Validate(func(s string) error {
								if s == "" {
									return fmt.Errorf("URL is required")
								}
								if !validator.IsValidURL(s) {
									return fmt.Errorf("invalid URL format")
								}
								return nil
							}),
					),
					huh.NewGroup(
						huh.NewConfirm().
							Title("Download video or audio?").
							Affirmative("Video").
							Negative("Audio").
							Value(&m.isVideo).
							Key("isVideo"),
					),
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("Select resolution").
							Options(
								huh.NewOption("720p", "720p"),
								huh.NewOption("1080p (recommended)", "1080p"),
								huh.NewOption("1440p", "1440p"),
								huh.NewOption("4K", "4K"),
							).
							Value(&m.resolution).
							Key("resolution"),
					).WithHideFunc(func() bool { return !m.isVideo }),
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("Select format").
							OptionsFunc(func() []huh.Option[string] {
								if m.isVideo {
									return []huh.Option[string]{
										huh.NewOption("MP4 (recommended)", "mp4"),
										huh.NewOption("MKV", "mkv"),
										huh.NewOption("WebM", "webm"),
									}
								}
								return []huh.Option[string]{
									huh.NewOption("MP3 (recommended)", "mp3"),
									huh.NewOption("M4A", "m4a"),
									huh.NewOption("FLAC", "flac"),
									huh.NewOption("WAV", "wav"),
								}
							}, &m.isVideo).
							Value(&m.format).
							Key("format"),
					),
					huh.NewGroup(
						huh.NewInput().
							Title("Download folder").
							Placeholder("~/Downloads").
							Value(&m.folder).
							Key("folder").
							Validate(func(s string) error {
								if s == "" {
									return fmt.Errorf("folder is required")
								}
								expanded := os.ExpandEnv(s)
								if _, err := os.Stat(expanded); os.IsNotExist(err) {
									return fmt.Errorf("folder does not exist")
								}
								return nil
							}),
					),
				)
				return m, m.form.Init()
			}
			return m, tea.Quit
		}
		if m.done && m.err != nil {
			if msg.String() == "q" {
				return m, tea.Quit
			}
			if msg.String() == "r" {
				m.done = false
				m.err = nil
				m.downloading = false
			}
			return m, nil
		}
		if m.downloading {
			return m, nil
		}
		if msg.String() == "d" {
			m.skipToDl = true
			m.downloading = true
			m.folder = os.ExpandEnv(m.folder)
			return m, tea.Batch(downloadCmd(m), tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return spinnerTick(t) }))
		}
		if m.form.State == huh.StateCompleted {
			if msg.String() == "enter" {
				m.downloading = true
				m.folder = os.ExpandEnv(m.folder)
				return m, tea.Batch(downloadCmd(m), tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return spinnerTick(t) }))
			}
		}
	case spinnerTick:
		m.spinnerFrame = (m.spinnerFrame + 1) % len(spinnerChars)
		return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return spinnerTick(t) })
	case downloadCompleteMsg:
		m.downloading = false
		m.done = true
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.successMsg = msg.title
			m.url = msg.url
		}
		return m, nil
	}

	form, cmd := m.form.Update(msg)
	m.form = form.(*huh.Form)
	return m, cmd
}

var spinnerChars = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func (m *model) View() tea.View {
	v := tea.NewView("")
	v.WindowTitle = "ytui"
	v.AltScreen = true

	if m.downloading {
		spinner := SpinnerStyle.Render(spinnerChars[m.spinnerFrame])
		v.SetContent(fmt.Sprintf("\n  %s Downloading...\n\n  Press Ctrl+C to cancel.\n", spinner))
		return v
	}

	if m.done && m.err != nil {
		errView := fmt.Sprintf("\n  %s\n\n  Error: %s\n\n  %s\n",
			TitleStyle.Render("Download Failed"),
			ErrorStyle.Render(m.err.Error()),
			"[ Press R to retry, Q to quit ]")
		v.SetContent(errView)
		return v
	}

	if m.done && m.successMsg != "" {
		successView := fmt.Sprintf("\n  %s\n\n  Downloaded: %s (%s)\n\n  %s\n",
			TitleStyle.Render("Download Complete"),
			SuccessStyle.Render(m.successMsg),
			SuccessStyle.Render(m.url),
			"[ Press R to download another, Q to quit ]")
		v.SetContent(successView)
		return v
	}

	if m.form.State == huh.StateCompleted {
		res := m.resolution
		f := m.format
		if !m.isVideo {
			res = "N/A"
			if f == "mp4" || f == "mkv" || f == "webm" {
				f = "mp3"
			}
		}

		confirmView := fmt.Sprintf("\n  %s\n\n  URL: %s\n  Type: %s\n  Resolution: %s\n  Format: %s\n  Folder: %s\n\n  %s\n",
			TitleStyle.Render("Confirm Download"),
			m.url,
			getTypeDisplay(m.isVideo),
			res,
			f,
			m.folder,
			"[ Press Enter to download ]")
		v.SetContent(confirmView)
		return v
	}

	view := m.form.View()
	view += "\n\n  Press D to skip prompts and use defaults"
	v.SetContent(view)
	return v
}

func getTypeDisplay(isVideo bool) string {
	if isVideo {
		return "Video"
	}
	return "Audio"
}

type spinnerTick time.Time

type downloadCompleteMsg struct {
	title string
	url   string
	err   error
}

func downloadCmd(m *model) tea.Cmd {
	return func() tea.Msg {
		// First get the video title
		title, err := ytdlp.GetVideoTitle(m.url)
		if err != nil {
			title = m.url // fallback to URL if we can't get title
		}

		opts := ytdlp.Options{
			URL:            m.url,
			IsAudio:        !m.isVideo,
			Resolution:     m.resolution,
			Format:         m.format,
			DownloadFolder: m.folder,
		}

		args, err := opts.BuildArgs()
		if err != nil {
			SendError(err.Error())
			return downloadCompleteMsg{err: err}
		}

		args = append([]string{"-q"}, args...)
		cmd := exec.Command("yt-dlp", args...)
		cmd.Stdout = nil
		cmd.Stderr = nil

		err = cmd.Run()
		if err != nil {
			SendError(err.Error())
			return downloadCompleteMsg{err: err}
		}

		SendSuccess(title)
		return downloadCompleteMsg{title: title, url: m.url, err: nil}
	}
}

func RunWithClipboard() error {
	cliURL := ""

	clipContent, err := clipboard.Read()
	if err == nil && validator.IsValidURL(clipContent) {
		cliURL = clipContent
	}

	m := New(cliURL)
	p := tea.NewProgram(m)

	_, err = p.Run()
	if err != nil {
		return err
	}

	if m.done && m.err == nil {
		os.Exit(0)
	}

	return nil
}
