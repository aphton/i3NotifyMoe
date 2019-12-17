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

	data, err := get(APIEndpointNickToUser+nm.configSvc.GetConfiguration().Username, &User{})
	if err != nil {
		return nil, err
	}
	user := data.(*User)

	data, err = get(APIEndpointAnimeList+user.UserID, &Animelist{})
	if err != nil {
		return nil, err
	}
	animelist := data.(*Animelist)

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
		data, err := get(APIEndpointEpisode+idNextEpisode, &Episode{})
		if err != nil {
			return nil, err
		}
		episode := data.(*Episode)

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

func get(endPoint string, dataType interface{}) (interface{}, error) {
	resp, err := http.Get(endPoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &dataType)
	return dataType, err
}

func animeByListItem(item *AnimelistItem, wg *sync.WaitGroup, animes chan *AnimeWithEpisode) {
	defer wg.Done()

	data, err := get(APIEndpointAnime+item.AnimeID, &Anime{})
	if err != nil {
		return
	}

	animeWithEpisode := AnimeWithEpisode{
		CurEpisode: item.Episodes,
		Anime:      *data.(*Anime),
	}

	animes <- &animeWithEpisode
}
