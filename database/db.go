package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"

	"porygon/config"
)

type GymStats struct {
	Valor       int
	Mystic      int
	Instinct    int
	Uncontested int
}

type PokeStats struct {
	Scanned      int
	Hundo        int
	Nundo        int
	Shiny        int
	ShinySpecies int
}

type RaidStats struct {
	Level int
	Raid  int
	Egg   int
}

type TypeCountStats struct {
	Type  int
	Count int
}

func DbConn(config config.Config) (*sqlx.DB, error) {
	DB, err := sqlx.Open("mysql", config.Database.Username+":"+config.Database.Password+"@tcp("+config.Database.Host+":"+config.Database.Port+")/"+config.Database.Name)
	return DB, err
}

func GetPokeStats(db *sqlx.DB) (PokeStats, error) {
	pokeStats := PokeStats{}

	err := db.Get(&pokeStats, `
		SELECT 
			(SELECT COALESCE(SUM(count), 0) FROM pokemon_stats WHERE date = CURDATE()) AS scanned,
			(SELECT COALESCE(SUM(count), 0) FROM pokemon_hundo_stats WHERE date = CURDATE()) AS hundo,
			(SELECT COALESCE(SUM(count), 0) FROM pokemon_nundo_stats WHERE date = CURDATE()) AS nundo,
			(SELECT COALESCE(SUM(count), 0) FROM pokemon_shiny_stats WHERE date = CURDATE()) AS shiny,
			(SELECT COUNT(DISTINCT pokemon_id) FROM pokemon_shiny_stats WHERE date = CURDATE()) AS shinyspecies
	`)
	if err != nil {
		return PokeStats{}, err
	}

	return pokeStats, nil
}

func GetRaidStats(db *sqlx.DB) ([]RaidStats, error) {
	var raidStatsList []RaidStats

	query := `
        SELECT 
            raid_level AS level,
            SUM(CASE WHEN raid_end_timestamp > UNIX_TIMESTAMP() THEN 1 ELSE 0 END) AS raid,
            SUM(CASE WHEN raid_battle_timestamp > UNIX_TIMESTAMP() AND (raid_spawn_timestamp IS NULL OR raid_spawn_timestamp <= UNIX_TIMESTAMP()) THEN 1 ELSE 0 END) AS egg
        FROM 
            gym
        GROUP BY 
            raid_level
        HAVING 
            raid > 0 OR egg > 0
    `
	err := db.Select(&raidStatsList, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return raidStatsList, nil
}

func GetGymStats(db *sqlx.DB) (GymStats, error) {
	gymStats := GymStats{}
	var err error
	err = db.Get(&gymStats, `
		SELECT 
            SUM(CASE WHEN team_id = 0 THEN 1 ELSE 0 END) AS uncontested,
            SUM(CASE WHEN team_id = 1 THEN 1 ELSE 0 END) AS valor,
            SUM(CASE WHEN team_id = 2 THEN 1 ELSE 0 END) AS mystic,
            SUM(CASE WHEN team_id = 3 THEN 1 ELSE 0 END) AS instinct
		FROM 
			gym 
		WHERE 
			updated > UNIX_TIMESTAMP() - 4 * 60 * 60
	`)

	if err != nil {
		log.Println(err)
		return GymStats{}, err
	}

	return gymStats, nil
}

func GetPokestopStats(db *sqlx.DB) (int, error) {
	var totalCount int
	err := db.QueryRow("SELECT COUNT(*) FROM pokestop WHERE quest_expiry > UNIX_TIMESTAMP()").Scan(&totalCount)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func GetRewardStats(db *sqlx.DB) ([]TypeCountStats, error) {
	var rewardStatsList []TypeCountStats

	err := db.Select(&rewardStatsList, `
		SELECT reward_type as type, SUM(count) AS count
		FROM (
			SELECT quest_reward_type AS reward_type, COUNT(*) AS count
			FROM pokestop
			WHERE quest_expiry > UNIX_TIMESTAMP() AND quest_reward_type IS NOT NULL
			GROUP BY quest_reward_type
			UNION ALL
			SELECT alternative_quest_reward_type AS reward_type, COUNT(*) AS count
			FROM pokestop
			WHERE quest_expiry > UNIX_TIMESTAMP() AND alternative_quest_reward_type IS NOT NULL
			GROUP BY alternative_quest_reward_type
		) AS subquery
		GROUP BY reward_type
    `)
	if err != nil {
		log.Println(err)
		return rewardStatsList, err
	}

	return rewardStatsList, nil
}

func GetLureStats(db *sqlx.DB) ([]TypeCountStats, error) {
	var lureStatsList []TypeCountStats

	err := db.Select(&lureStatsList, `
		SELECT 
		    lure_id as type, COUNT(*) as count
		FROM 
		    pokestop
		WHERE
		    lure_expire_timestamp > UNIX_TIMESTAMP()
		GROUP BY lure_id
    `)
	if err != nil {
		return lureStatsList, err
	}

	return lureStatsList, nil
}

func GetRocketStats(db *sqlx.DB) ([]TypeCountStats, error) {
	var rocketStatsList []TypeCountStats

	err := db.Select(&rocketStatsList, `
		SELECT 
		    display_type as type, COUNT(*) as count
		FROM
		    incident
		WHERE
		    expiration > UNIX_TIMESTAMP() AND display_type < 8
		GROUP BY display_type
    `)
	if err != nil {
		return rocketStatsList, err
	}

	return rocketStatsList, nil
}

func GetEventStats(db *sqlx.DB) ([]TypeCountStats, error) {
	var eventStatsList []TypeCountStats

	err := db.Select(&eventStatsList, `
		SELECT 
		    display_type as type, COUNT(*) as count
		FROM
		    incident
		WHERE
		    expiration > UNIX_TIMESTAMP() AND display_type >= 8
		GROUP BY display_type
    `)
	if err != nil {
		return eventStatsList, err
	}

	return eventStatsList, nil
}

func GetRoutesStats(db *sqlx.DB) (int, error) {
	var totalCount int
	err := db.QueryRow("SELECT COUNT(*) FROM route WHERE type = 1").Scan(&totalCount)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}
