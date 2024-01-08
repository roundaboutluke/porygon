package config

type Config struct {
	Discord     discord
	Coordinates coordinates
	Database    database
	API         api
	Config      config
}

type config struct {
	RefreshInterval     int
	IncludeActiveCounts bool
	EmbedTitle          string
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
	Emojis     emojistruct
}

type emojistruct struct {
	Valor       string
	Mystic      string
	Instinct    string
	Uncontested string
	Pokestop    string
	Normal      string
	Glacial     string
	Mossy       string
	Magnetic    string
	Rainy       string
	Sparkly     string
	Scanned     string
	Hundo       string
	Nundo       string
	Shinies     string
	Grunt       string
	Leader      string
	Giovanni    string
	Kecleon     string
	Showcase    string
	Route       string
	Level1      string
	Level3      string
	Level4      string
	Level5      string
	Mega        string
	Elite       string
	Items       string
	Encounter   string
	Stardust    string
	MegaEnergy  string
}
