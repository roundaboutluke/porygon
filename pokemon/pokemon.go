package pokemon

import "strings"

type Incident struct {
	ID    int
	Emoji string
}

type Raid struct {
	ID    int
	Emoji string
}

type Reward struct {
	ID    int
	Emoji string
}

type Pokestop struct {
	ID    int
	Emoji string
}

type Team struct {
	ID    int
	Emoji string
}

type Lure struct {
	ID    int
	Emoji string
}

// honestly don't think this belongs here XD but will leave it until i find it a home
func FormatEmoji(emoji string) string {
	if strings.Contains(emoji, "<") && strings.Contains(emoji, ">") {
		return emoji
	} else if strings.Contains(emoji, ":") {
		return "<" + emoji + ">"
	}
	return emoji
}
