package notifymoeweb

import "github.com/Aphton/i3NotifyMoe/config"

// User notify moe user
type User struct {
	Nick   string
	UserID string
}

// Animelist notify moe Animelist
type Animelist struct {
	Items []AnimelistItem
}

// AnimelistItem notify moe AnimelistItem
type AnimelistItem struct {
	AnimeID  string
	Status   string
	Episodes int
}

// AnimeWithEpisode custom object wrapping/tracking episode for an Anime item
type AnimeWithEpisode struct {
	CurEpisode int
	Anime      Anime
}

// Anime notify moe Anime
type Anime struct {
	Title    AnimeTitle
	Status   string
	Episode  int
	Episodes []string
}

// AnimeTitle notify moe AnimeTitle
type AnimeTitle struct {
	Romaji string
}

// Episode notify moe Episode
type Episode struct {
	Number     int
	AiringDate AiringDate
}

// AiringDate notify moe AiringDate
type AiringDate struct {
	Start string
	End   string
}

type notifyMoe struct {
	configSvc config.ConfigurationService
}
