package config

type Config struct {
	Discord     discord
	Coordinates coordinates
	Database    database
	API         api
	Config      config
	LevelEmoji  map[string]string `toml:"level_emoji"`
	RewardEmoji map[string]string `toml:"reward_emoji"`
	LureEmoji   map[string]string `toml:"lure_emoji"`
	RocketEmoji map[string]string `toml:"rocket_emoji"`
	EventEmoji  map[string]string `toml:"event_emoji"`
}

type config struct {
	RefreshInterval      int
	ErrorRefreshInterval int
	IncludeActiveCounts  bool
	EmbedTitle           string
}

type api struct {
	URL    string
	Secret string
}
type database struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
}

type coordinates struct {
	Min struct {
		Latitude  float64
		Longitude float64
	}
	Max struct {
		Latitude  float64
		Longitude float64
	}
}

type discord struct {
	Token      string
	ChannelIDs []string
}
