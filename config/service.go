package config

import "io"

// ConfigurationService interface for configuration
type ConfigurationService interface {
	LoadPersistedState(reader io.Reader) error
	PersistState(writer io.Writer) error
	SetStateItems(stateItems []StateItem) error
	GetConfiguration() *Configuration
	GetState() *State
}
