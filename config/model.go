package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Context application state
type context struct {
	config Configuration
	state  State
}

// Configuration contains environment variables that are set and influence the behavior of the application
type Configuration struct {
	Username    string
	BlockButton string
}

// State contains a list of each category of queried animes (released, today, unreleasd)
type State struct {
	ReleasedAnimes   []StateItem
	TodayAnimes      []StateItem
	UnreleasedAnimes []StateItem
}

// StateItem describes an anime
type StateItem struct {
	Title    string
	Episode  int
	DateTime string
}

// IsRightButtonClicked checks whether the environment variable for block button indicates a right button click
func (cfg Configuration) IsRightButtonClicked() bool {
	return cfg.BlockButton == BlockButtonRightClick
}

// FromString converts a formatted StateItem string back (deserialization)
func (si *StateItem) FromString(str string) error {
	items := strings.Split(str, ";")
	if len(items) != 3 {
		return errors.New("error")
	}
	ep, err := strconv.ParseInt(items[1], 10, 64)
	si.Title = items[0]
	if err == nil {
		si.Episode = int(ep)
	}
	si.DateTime = items[2]
	return nil
}

// DatetimeToDate converts datetime string to a timezone conscious datetime object and extracts the date
func (si StateItem) DatetimeToDate(location *time.Location) string {
	datetime, _ := time.Parse(time.RFC3339, si.DateTime)
	datetime = datetime.In(location)
	return datetime.Format("2006-01-02")
}

// DatetimeToTime converts datetime string to a timezone conscious datetime object and extracts the time
func (si StateItem) DatetimeToTime(location *time.Location) string {
	datetime, _ := time.Parse(time.RFC3339, si.DateTime)
	datetime = datetime.In(location)
	return datetime.Format("15:04:05")
}

func (si StateItem) String() string {
	return fmt.Sprintf("%s;%d;%s", si.Title, si.Episode, si.DateTime)
}

func (s *State) String() string {
	return fmt.Sprintf("\uf008 %d/%d",
		len(s.ReleasedAnimes)+len(s.TodayAnimes),
		len(s.ReleasedAnimes)+len(s.TodayAnimes)+len(s.UnreleasedAnimes))
}
