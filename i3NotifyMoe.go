package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Aphton/i3NotifyMoe/config"
	"github.com/Aphton/i3NotifyMoe/desktopnotify"
	"github.com/Aphton/i3NotifyMoe/notifymoeweb"
)

func getStateFilename() string {
	result, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Join(filepath.Dir(result), config.StateFilename)
}

func loadAndNotify(configSvc config.ConfigurationService, notifySvc desktopnotify.NotifyService) error {
	loc := time.Now().Location()
	var released, today, unreleased []string
	for _, r := range configSvc.GetState().ReleasedAnimes {
		released = append(released, fmt.Sprintf("%s - %d", r.Title, r.Episode))
	}
	for _, t := range configSvc.GetState().TodayAnimes {
		today = append(today, fmt.Sprintf("%s %s - %d", t.DatetimeToTime(loc), t.Title, t.Episode))
	}
	for _, u := range configSvc.GetState().UnreleasedAnimes {
		unreleased = append(unreleased, fmt.Sprintf("%s %s - %d", u.DatetimeToDate(loc), u.Title, u.Episode))
	}
	if err := notifySvc.NotifySend(released, today, unreleased); err != nil {
		return err
	}

	fmt.Println(configSvc.GetState())
	return nil
}

func fetchAndDisplay(notifyMoeSvc notifymoeweb.NotifyMoeService, configSvc config.ConfigurationService) error {
	animeStates, err := notifyMoeSvc.FetchCurrentlyWatchingAndAiringAnimes()
	if err != nil {
		return err
	}
	if err := configSvc.SetStateItems(animeStates); err != nil {
		return err
	}
	file, err := os.Create(getStateFilename())
	if err != nil {
		return err
	}
	defer file.Close()
	if err := configSvc.PersistState(file); err != nil {
		return err
	}
	fmt.Println(configSvc.GetState())
	return nil
}

func run() error {
	configSvc := config.NewConfigurationService()
	if configSvc.GetConfiguration().IsRightButtonClicked() {
		file, err := os.Open(getStateFilename())
		if err != nil {
			return err
		}
		err = configSvc.LoadPersistedState(file)
		if err != nil {
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}
	}

	notifyMoeSvc := notifymoeweb.NewNotifyMoeService(configSvc)

	notifySvc, err := desktopnotify.NewNotifyService(desktopnotify.JadeTemplate)
	if err != nil {
		return err
	}

	if configSvc.GetConfiguration().IsRightButtonClicked() {
		return loadAndNotify(configSvc, notifySvc)
	}

	return fetchAndDisplay(notifyMoeSvc, configSvc)
}

func main() {
	if err := run(); err != nil {
		log.Println(err)
		fmt.Println("error")
	}
}
