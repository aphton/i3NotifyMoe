package desktopnotify

// NotifyService defines a method to send to "notify-send"
type NotifyService interface {
	NotifySend(released, today, unreleased []string) error
}
