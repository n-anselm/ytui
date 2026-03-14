package clipboard

import (
	"os/exec"
	"strings"
)

func Read() (string, error) {
	tools := []string{"wl-paste", "xclip", "xsel"}

	for _, tool := range tools {
		var cmd *exec.Cmd
		switch tool {
		case "wl-paste":
			cmd = exec.Command("wl-paste")
		case "xclip":
			cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
		case "xsel":
			cmd = exec.Command("xsel", "--clipboard", "--output")
		}

		out, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(out)), nil
		}
	}

	return "", nil
}
