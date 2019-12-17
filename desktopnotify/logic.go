package desktopnotify

import (
	"os/exec"

	"github.com/go-floki/jade"
)

// NewNotifyService creates a new NotifyService
func NewNotifyService(templateStr string) (NotifyService, error) {
	var ns notifyService
	tpl, err := jade.Compile(templateStr, jade.Options{})
	ns.tpl = tpl
	return &ns, err
}

// NotifySend prepares and launches a "notify-send" command
func (ns *notifyService) NotifySend(released, today, unreleased []string) error {
	if err := ns.tpl.Execute(ns, map[string]interface{}{
		"Released":   released,
		"Today":      today,
		"Unreleased": unreleased,
	}); err != nil {
		return err
	}

	cmd := exec.Command("notify-send", "notify.moe", ns.buffer.String())
	return cmd.Run()
}
func (ns *notifyService) Write(p []byte) (n int, err error) {
	ns.buffer.WriteString(string(p))
	return len(p), nil
}
