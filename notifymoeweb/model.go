package notifymoeweb

import "github.com/Aphton/i3NotifyMoe/config"

// User from notify moe api
type User struct {
	Nick   string
	UserID string
}

// Animelist from notify moe api
type Animelist struct {
	Items []AnimelistItem
}

// AnimelistItem from notify moe api
type AnimelistItem struct {
	AnimeID  string
	Status   string
	Episodes int
}

// AnimeWithEpisode custom helper struct, wrapping episode around an Anime struct
type AnimeWithEpisode struct {
	CurEpisode int
	Anime      Anime
}

// Anime from notify moe api
type Anime struct {
	Title    AnimeTitle
	Status   string
	Episode  int
	Episodes []string
}

// AnimeTitle from notify moe api
type AnimeTitle struct {
	Romaji string
}

// Episode from notify moe api
type Episode struct {
	Number     int
	AiringDate AiringDate
}

// AiringDate from notify moe api
type AiringDate struct {
	Start string
	End   string
}

type notifyMoe struct {
	configSvc config.ConfigurationService
}
