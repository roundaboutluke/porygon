package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"

	"Porygon/config"
	"Porygon/pokemon"
)

var DB *sql.DB

func DbConn(config config.Config) (*sql.DB, error) {
	DB, err := sql.Open("mysql", config.Database.Username+":"+config.Database.Password+"@tcp("+config.Database.Host+":"+config.Database.Port+")/"+config.Database.Name)
	return DB, err
}

func RaidStats(db *sql.DB, config config.Config) (string, error) {
	raidEggStats := ""
	var err error

	raids := []pokemon.Raid{
		{ID: 1, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Level1)},
		{ID: 3, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Level3)},
		{ID: 4, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Level4)},
		{ID: 5, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Level5)},
		{ID: 6, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Mega)},
		{ID: 9, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Elite)},
	}

	for _, raid := range raids {
		var activeRaidCount, activeEggCount int
		if raid.ID == 5 {
			err = db.QueryRow("SELECT COUNT(*) FROM gym WHERE raid_level IN (5, 8) AND raid_battle_timestamp <= UNIX_TIMESTAMP() AND raid_end_timestamp > UNIX_TIMESTAMP()").Scan(&activeRaidCount)
			err = db.QueryRow("SELECT COUNT(*) FROM gym WHERE raid_level IN (5, 8) AND raid_battle_timestamp > UNIX_TIMESTAMP() AND (raid_spawn_timestamp IS NULL OR raid_spawn_timestamp <= UNIX_TIMESTAMP())").Scan(&activeEggCount)
		} else if raid.ID == 6 {
			err = db.QueryRow("SELECT COUNT(*) FROM gym WHERE raid_level IN (6, 7, 10) AND raid_battle_timestamp <= UNIX_TIMESTAMP() AND raid_end_timestamp > UNIX_TIMESTAMP()").Scan(&activeRaidCount)
			err = db.QueryRow("SELECT COUNT(*) FROM gym WHERE raid_level IN (6, 7, 10) AND raid_battle_timestamp > UNIX_TIMESTAMP() AND (raid_spawn_timestamp IS NULL OR raid_spawn_timestamp <= UNIX_TIMESTAMP())").Scan(&activeEggCount)
		} else {
			err = db.QueryRow("SELECT COUNT(*) FROM gym WHERE raid_level = ? AND raid_battle_timestamp <= UNIX_TIMESTAMP() AND raid_end_timestamp > UNIX_TIMESTAMP()", raid.ID).Scan(&activeRaidCount)
			err = db.QueryRow("SELECT COUNT(*) FROM gym WHERE raid_level = ? AND raid_battle_timestamp > UNIX_TIMESTAMP() AND (raid_spawn_timestamp IS NULL OR raid_spawn_timestamp <= UNIX_TIMESTAMP())", raid.ID).Scan(&activeEggCount)
		}

		if err != nil {
			return raidEggStats, err
		}
		if activeRaidCount > 0 || activeEggCount > 0 {
			raidEggStats += fmt.Sprintf("%s Hatched: %d | Eggs: %d\n", raid.Emoji, activeRaidCount, activeEggCount)

		}
	}
	return raidEggStats, err

}

func GymStats(db *sql.DB, config config.Config) (string, error) {
	gymStats := ""
	var err error

	teams := []pokemon.Team{
		{ID: 1, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Valor)},
		{ID: 2, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Mystic)},
		{ID: 3, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Instinct)},
		{ID: 0, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Uncontested)},
	}

	for _, team := range teams {
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM gym WHERE team_id = ? AND updated > UNIX_TIMESTAMP() - 4 * 60 * 60", team.ID).Scan(&count)

		// not sure behaviour if error, is Sprintf count still okay? temp add error nil
		if err != nil {
			return gymStats, err
		}

		gymStats += fmt.Sprintf("%s %d ", team.Emoji, count)
	}
	return gymStats, err
}

func PokestopStats(db *sql.DB, config config.Config) (string, error) {
	pokestopStats := ""
	var err error

	pokestops := []pokemon.Pokestop{
		{ID: 1, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Pokestop)},
		// Add more if needed
	}

	for _, pokestop := range pokestops {
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM pokestop WHERE quest_expiry > UNIX_TIMESTAMP()").Scan(&count)

		if err != nil {
			return pokestopStats, err
		}

		pokestopStats += fmt.Sprintf("%s %d ", pokestop.Emoji, count)
	}

	return pokestopStats, err
}

func RewardStats(db *sql.DB, config config.Config) (string, error) {
	rewardStats := ""
	var err1, err2 error

	rewards := []pokemon.Reward{
		{ID: 2, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Items)},
		{ID: 7, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Encounter)},
		{ID: 3, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Stardust)},
		{ID: 12, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.MegaEnergy)},
	}
	for _, reward := range rewards {
		var count1, count2 int
		err1 = db.QueryRow("SELECT COUNT(*) FROM pokestop WHERE quest_reward_type = ? AND quest_expiry > UNIX_TIMESTAMP()", reward.ID).Scan(&count1)
		err2 = db.QueryRow("SELECT COUNT(*) FROM pokestop WHERE alternative_quest_reward_type = ? AND quest_expiry > UNIX_TIMESTAMP()", reward.ID).Scan(&count2)

		if err1 != nil {
			return rewardStats, err1
		}

		if err2 != nil {
			return rewardStats, err2
		}

		count := count1 + count2
		rewardStats += fmt.Sprintf("%s %d ", reward.Emoji, count)
	}
	return rewardStats, nil
}

func LureStats(db *sql.DB, config config.Config) (string, error) {
	lureStats := ""
	var err error

	lures := []pokemon.Lure{
		{ID: 501, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Normal)},
		{ID: 502, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Glacial)},
		{ID: 503, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Mossy)},
		{ID: 504, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Magnetic)},
		{ID: 505, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Rainy)},
		{ID: 506, Emoji: pokemon.FormatEmoji(config.Discord.Emojis.Sparkly)},
	}

	for _, lure := range lures {
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM pokestop WHERE lure_id = ? AND lure_expire_timestamp > UNIX_TIMESTAMP()", lure.ID).Scan(&count)
		if err != nil {
			return lureStats, err
		}

		lureStats += fmt.Sprintf("%s %d ", lure.Emoji, count)
	}

	return lureStats, err
}
