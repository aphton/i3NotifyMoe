package notifymoeweb

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/Aphton/i3NotifyMoe/config"
)

// NewNotifyMoeService creates a new NotifyMoeService
func NewNotifyMoeService(configSvc config.ConfigurationService) NotifyMoeService {
	return &notifyMoe{
		configSvc: configSvc,
	}
}

// FetchCurrentlyWatchingAndAiringAnimes crawles the API endpoints for currently watching & airing animes, to track their state
func (nm *notifyMoe) FetchCurrentlyWatchingAndAiringAnimes() ([]config.StateItem, error) {
	if len(nm.configSvc.GetConfiguration().Username) == 0 {
		return nil, errors.New("Environment var '" + config.EnvUsername + "' not set!")
	}

	user, err := nickToUser(nm.configSvc.GetConfiguration())
	if err != nil {
		return nil, err
	}

	animelist, err := animelistOfUser(user)
	if err != nil {
		return nil, err
	}

	var watchingAnimes []AnimelistItem
	for _, item := range animelist.Items {
		if item.Status != AnimeListItemStatusWatching {
			continue
		}
		watchingAnimes = append(watchingAnimes, item)
	}

	var wg sync.WaitGroup
	animesWithEpisode := make(chan *AnimeWithEpisode, len(watchingAnimes))
	for idx := range watchingAnimes {
		wg.Add(1)
		go animeByListItem(&watchingAnimes[idx], &wg, animesWithEpisode)
	}
	wg.Wait()
	close(animesWithEpisode)

	var result []config.StateItem
	for animeWithEpisode := range animesWithEpisode {
		if animeWithEpisode.Anime.Status != AnimeStatusCurrent {
			continue
		}
		if animeWithEpisode.CurEpisode >= len(animeWithEpisode.Anime.Episodes) {
			continue // ERROR! should never happen
		}

		idNextEpisode := animeWithEpisode.Anime.Episodes[animeWithEpisode.CurEpisode]
		episode, err := episdeByID(idNextEpisode)
		if err != nil {
			continue
		}

		result = append(result, config.StateItem{
			Title:    animeWithEpisode.Anime.Title.Romaji,
			Episode:  episode.Number,
			DateTime: episode.AiringDate.End,
		})
	}

	return result, nil
}

func get(endPoint string) (*string, error) {
	resp, err := http.Get(endPoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	b := string(body)
	return &b, err
}

func nickToUser(cfg *config.Configuration) (*User, error) {
	body, err := get(APIEndpointNickToUser + cfg.Username)
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal([]byte(*body), &user)
	return &user, err
}

func animelistOfUser(user *User) (*Animelist, error) {
	body, err := get(APIEndpointAnimeList + user.UserID)
	if err != nil {
		return nil, err
	}

	var animelist Animelist
	err = json.Unmarshal([]byte(*body), &animelist)
	return &animelist, err
}

func episdeByID(episodeID string) (*Episode, error) {
	body, err := get(APIEndpointEpisode + episodeID)
	if err != nil {
		return nil, err
	}

	var episode Episode
	err = json.Unmarshal([]byte(*body), &episode)
	return &episode, err
}

func animeByListItem(item *AnimelistItem, wg *sync.WaitGroup, animes chan *AnimeWithEpisode) {
	defer wg.Done()

	body, err := get(APIEndpointAnime + item.AnimeID)
	if err != nil {
		return
	}

	animeWithEpisode := AnimeWithEpisode{}
	err = json.Unmarshal([]byte(*body), &animeWithEpisode.Anime)
	if err != nil {
		return
	}
	animeWithEpisode.CurEpisode = item.Episodes

	animes <- &animeWithEpisode
}
