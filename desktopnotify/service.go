package desktopnotify

// NotifyService defines a method to execute a "notify-send" command with a jade-templated argument
type NotifyService interface {
	NotifySend(released, today, unreleased []string) error
}
