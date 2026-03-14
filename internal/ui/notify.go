package ui

import (
	"os/exec"
	"strings"
)

func SendNotification(title, message string) error {
	cmd := exec.Command("notify-send", "-a", "ytui", title, message)
	return cmd.Run()
}

func SendSuccess(filename string) error {
	return SendNotification("Download Complete", filename+" has been downloaded successfully.")
}

func SendError(errMsg string) error {
	return SendNotification("Download Failed", errMsg)
}

func SendInfo(message string) error {
	return SendNotification("ytui", message)
}

func NotifyFallback(original, actual string) error {
	msg := strings.TrimSpace(original) + " not available. Downloaded " + strings.TrimSpace(actual) + " instead."
	return SendNotification("Resolution Fallback", msg)
}
