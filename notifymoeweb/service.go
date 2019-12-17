package notifymoeweb

import "github.com/Aphton/i3NotifyMoe/config"

// NotifyMoeService api
type NotifyMoeService interface {
	FetchCurrentlyWatchingAndAiringAnimes() ([]config.StateItem, error)
}
