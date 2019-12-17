package config

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

// NewConfigurationService creates a new ConfigurationService
func NewConfigurationService() ConfigurationService {
	return &context{
		config: loadConfig(),
	}
}

// LoadPersistedState loads state from the reader (file)
func (ctx *context) LoadPersistedState(reader io.Reader) error {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	var stateItems []StateItem
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		var si StateItem
		err := si.FromString(line)
		if err != nil {
			continue
		}
		stateItems = append(stateItems, si)
	}
	return ctx.SetStateItems(stateItems)
}

// PersistState sorts the lists in ctx.state and then persists them
func (ctx *context) PersistState(writer io.Writer) error {
	if err := sortState(ctx); err != nil {
		return err
	}
	var out bytes.Buffer
	for _, r := range ctx.state.ReleasedAnimes {
		out.WriteString(r.String())
		out.WriteString("\n")
	}
	for _, t := range ctx.state.TodayAnimes {
		out.WriteString(t.String())
		out.WriteString("\n")
	}
	for _, u := range ctx.state.UnreleasedAnimes {
		out.WriteString(u.String())
		out.WriteString("\n")
	}
	_, err := writer.Write(out.Bytes())
	return err
}

// SetStateItems loads & populates states from an unsorted list of StateItems
func (ctx *context) SetStateItems(stateItems []StateItem) error {
	now := time.Now()
	nowDateInt := timeToDateInt(&now)

	ctx.state.ReleasedAnimes = nil
	ctx.state.TodayAnimes = nil
	ctx.state.UnreleasedAnimes = nil

	for _, anime := range stateItems {
		endDate, err := time.Parse(time.RFC3339, anime.DateTime)
		if err != nil {
			return err
		}
		endDate = endDate.In(now.Location())
		endDateInt := timeToDateInt(&endDate)

		if endDateInt <= nowDateInt {
			if endDateInt == nowDateInt && now.Before(endDate) {
				ctx.state.TodayAnimes = append(ctx.state.TodayAnimes, anime)
			} else {
				ctx.state.ReleasedAnimes = append(ctx.state.ReleasedAnimes, anime)
			}
		} else {
			ctx.state.UnreleasedAnimes = append(ctx.state.UnreleasedAnimes, anime)
		}
	}
	return nil
}

func (ctx *context) GetConfiguration() *Configuration {
	return &ctx.config
}

func (ctx *context) GetState() *State {
	return &ctx.state
}

func timeToDateInt(t *time.Time) int {
	y, m, d := t.Date()
	return (y*100+int(m))*100 + d
}

func sortState(ctx *context) error {
	sorter := func(a, b StateItem) bool {
		endDate1, err := time.Parse(time.RFC3339, a.DateTime)
		if err != nil {
			return false
		}
		endDate2, err := time.Parse(time.RFC3339, b.DateTime)
		if err != nil {
			return false
		}
		return endDate1.Before(endDate2)

	}
	sort.Slice(ctx.state.ReleasedAnimes, func(i, j int) bool {
		return sorter(ctx.state.ReleasedAnimes[i], ctx.state.ReleasedAnimes[j])
	})
	sort.Slice(ctx.state.TodayAnimes, func(i, j int) bool {
		return sorter(ctx.state.TodayAnimes[i], ctx.state.TodayAnimes[j])
	})
	sort.Slice(ctx.state.UnreleasedAnimes, func(i, j int) bool {
		return sorter(ctx.state.UnreleasedAnimes[i], ctx.state.UnreleasedAnimes[j])
	})
	return nil
}

func loadConfig() Configuration {
	return Configuration{
		Username:    os.Getenv(EnvUsername),
		BlockButton: os.Getenv(BlockButton),
	}
}
