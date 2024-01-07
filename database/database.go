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

func RaidStats(db *sql.DB, raids []pokemon.Raid) (string, error) {
	raidEggStats := ""
	var err error

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

		if activeRaidCount > 0 || activeEggCount > 0 {
			raidEggStats += fmt.Sprintf("%s Hatched: %d | Eggs: %d\n", raid.Emoji, activeRaidCount, activeEggCount)

		}
	}
	return raidEggStats, err

}

func GymStats(db *sql.DB, teams []pokemon.Team) (string, error) {
	gymStats := ""
	var err error

	for _, team := range teams {
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM gym WHERE team_id = ? AND updated > UNIX_TIMESTAMP() - 4 * 60 * 60", team.ID).Scan(&count)
		if err != nil {
			fmt.Println("error querying MariaDB,", err)
		}

		gymStats += fmt.Sprintf("%s %d ", team.Emoji, count)
	}
	return gymStats, err
}
