package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// WoW race to faction mapping
var raceToFaction = map[int]string{
	1:  "Alliance", // Human
	2:  "Horde",    // Orc
	3:  "Alliance", // Dwarf
	4:  "Alliance", // Night Elf
	5:  "Horde",    // Undead
	6:  "Horde",    // Tauren
	7:  "Alliance", // Gnome
	8:  "Horde",    // Troll
	9:  "Horde",    // Goblin
	10: "Horde",    // Blood Elf
	11: "Alliance", // Draenei
	22: "Horde",    // Worgen
	24: "Alliance", // Pandaren (Neutral)
	25: "Alliance", // Pandaren (Alliance)
	26: "Horde",    // Pandaren (Horde)
}

// Metrics
var (
	playersOnline = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_players_online",
			Help: "Number of players currently online",
		},
		[]string{"faction"},
	)

	playersTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_players_total",
			Help: "Total number of players",
		},
		[]string{"faction"},
	)

	playersByLevel = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_players_by_level",
			Help: "Number of players by level",
		},
		[]string{"level", "faction"},
	)

	playersByClass = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_players_by_class",
			Help: "Number of players by class",
		},
		[]string{"class", "faction"},
	)

	mailTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_mail_total",
			Help: "Total number of mail messages",
		},
	)

	mailByFaction = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_mail_by_faction",
			Help: "Number of mail messages by faction",
		},
		[]string{"faction"},
	)

	mailWithItems = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_mail_with_items",
			Help: "Number of mail messages with items",
		},
	)

	accountsTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_accounts_total",
			Help: "Total number of accounts",
		},
	)

	accountsOnline = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_accounts_online",
			Help: "Number of accounts currently online",
		},
	)

	accountsBanned = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_accounts_banned",
			Help: "Number of banned accounts",
		},
	)

	serverUptime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_server_uptime_seconds",
			Help: "Server uptime in seconds",
		},
	)

	serverMaxPlayers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_server_max_players",
			Help: "Maximum number of players recorded",
		},
	)

	auctionCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_auction_count",
			Help: "Number of active auctions by house (faction)",
		},
		[]string{"house"},
	)

	guildCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_guild_count",
			Help: "Number of guilds",
		},
	)

	maxLevelCharCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_max_level_characters",
			Help: "Number of max-level characters by faction",
		},
		[]string{"faction"},
	)

	unreadMailCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_unread_mail_count",
			Help: "Number of unread mail messages",
		},
	)

	gmAccountCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_gm_account_count",
			Help: "Number of accounts with GM level",
		},
	)

	lastServerRestart = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_server_last_restart_timestamp",
			Help: "Timestamp of the last server restart (unix time)",
		},
	)

	bannedCharCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_banned_characters",
			Help: "Number of banned characters",
		},
	)

	// Additional chat and activity metrics
	channelCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_channel_count",
			Help: "Number of chat channels",
		},
	)

	channelBans = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_channel_bans",
			Help: "Number of channel bans",
		},
	)

	logCountByType = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_log_count",
			Help: "Number of log entries by type",
		},
		[]string{"type"},
	)

	guildEventCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_guild_events",
			Help: "Number of guild events",
		},
	)

	moneyLogCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_money_logs",
			Help: "Number of money transaction logs",
		},
	)

	encounterLogCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_encounter_logs",
			Help: "Number of encounter logs",
		},
	)

	arenaLogCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_arena_logs",
			Help: "Number of arena fight logs",
		},
	)

	ipActionLogCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_ip_action_logs",
			Help: "Number of IP action logs",
		},
	)

	// Additional instance and raid metrics
	activeInstanceCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_active_instances",
			Help: "Number of active instances",
		},
	)

	instancesByDifficulty = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_instances_by_difficulty",
			Help: "Number of instances by difficulty",
		},
		[]string{"difficulty"},
	)

	completedEncounters = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_completed_encounters",
			Help: "Number of completed encounters by instance",
		},
		[]string{"instance_id"},
	)

	instanceResets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_instance_resets",
			Help: "Instance reset times by map and difficulty",
		},
		[]string{"map_id", "difficulty"},
	)

	charactersInInstances = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_characters_in_instances",
			Help: "Number of characters currently in instances",
		},
	)

	lfgDataCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_lfg_data",
			Help: "Number of LFG entries by state",
		},
		[]string{"state"},
	)

	lagReportsCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_lag_reports",
			Help: "Number of lag reports",
		},
	)

	instanceSavesCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_instance_saves",
			Help: "Number of saved instance states",
		},
	)

	// Additional networking metrics
	playerLatencyStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_player_latency",
			Help: "Player latency statistics",
		},
		[]string{"stat"},
	)

	ipBannedCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_ip_banned_count",
			Help: "Number of banned IP addresses",
		},
	)

	ipActionLogsByType = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_ip_action_logs_by_type",
			Help: "Number of IP action logs by type",
		},
		[]string{"type"},
	)

	lagReportsByType = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_lag_reports_by_type",
			Help: "Number of lag reports by type",
		},
		[]string{"lag_type"},
	)

	averageLatency = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_average_latency_ms",
			Help: "Average player latency in milliseconds",
		},
	)

	highLatencyPlayers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_high_latency_players",
			Help: "Number of players with high latency (>200ms)",
		},
	)

	networkActivityByIP = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_network_activity_by_ip",
			Help: "Network activity by IP address (top 10)",
		},
		[]string{"ip"},
	)

	// Battleground metrics
	battlegroundDeserters = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_battleground_deserters",
			Help: "Number of battleground deserters",
		},
	)

	battlegroundDesertersByType = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battleground_deserters_by_type",
			Help: "Number of battleground deserters by type",
		},
		[]string{"desertion_type"},
	)

	randomBattlegroundQueue = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_random_battleground_queue",
			Help: "Number of players in random battleground queue",
		},
	)

	battlegroundStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battleground_stats",
			Help: "Battleground statistics",
		},
		[]string{"stat"},
	)

	battlegroundsByType = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battlegrounds_by_type",
			Help: "Number of battlegrounds by type",
		},
		[]string{"battleground_type"},
	)

	battlegroundsByBracket = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battlegrounds_by_bracket",
			Help: "Number of battlegrounds by bracket",
		},
		[]string{"bracket"},
	)

	battlegroundWinsByFaction = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battleground_wins_by_faction",
			Help: "Number of battleground wins by faction",
		},
		[]string{"faction"},
	)

	battlegroundPlayerStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battleground_player_stats",
			Help: "Battleground player statistics",
		},
		[]string{"stat"},
	)

	battlegroundTemplates = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battleground_templates",
			Help: "Battleground template information",
		},
		[]string{"template_id", "script_name"},
	)

	battlegroundTemplateDetails = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battleground_template_details",
			Help: "Detailed battleground template information",
		},
		[]string{"template_id", "name", "min_level", "max_level", "min_players", "max_players"},
	)

	recentBattlegrounds = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_recent_battlegrounds",
			Help: "Recent battleground activity",
		},
		[]string{"time_period"},
	)
)

// WoW class names
var classNames = map[int]string{
	1:  "Warrior",
	2:  "Paladin",
	3:  "Hunter",
	4:  "Rogue",
	5:  "Priest",
	6:  "Death Knight",
	7:  "Shaman",
	8:  "Mage",
	9:  "Warlock",
	10: "Monk",
	11: "Druid",
	12: "Demon Hunter",
}

type Exporter struct {
	charactersDB *sql.DB
	authDB       *sql.DB
	worldDB      *sql.DB
}

func NewExporter(dsn string) (*Exporter, error) {
	// Create connection to characters database
	charactersDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Test the characters connection
	if err := charactersDB.Ping(); err != nil {
		return nil, err
	}

	// Create connection to auth database by modifying the DSN
	authDSN := strings.Replace(dsn, "/acore_characters?", "/acore_auth?", 1)
	authDB, err := sql.Open("mysql", authDSN)
	if err != nil {
		charactersDB.Close()
		return nil, err
	}

	// Test the auth connection
	if err := authDB.Ping(); err != nil {
		charactersDB.Close()
		authDB.Close()
		return nil, err
	}

	// Create connection to world database by modifying the DSN
	worldDSN := strings.Replace(dsn, "/acore_characters?", "/acore_world?", 1)
	worldDB, err := sql.Open("mysql", worldDSN)
	if err != nil {
		charactersDB.Close()
		authDB.Close()
		return nil, err
	}

	// Test the world connection
	if err := worldDB.Ping(); err != nil {
		charactersDB.Close()
		authDB.Close()
		worldDB.Close()
		return nil, err
	}

	return &Exporter{
		charactersDB: charactersDB,
		authDB:       authDB,
		worldDB:      worldDB,
	}, nil
}

func (e *Exporter) Close() {
	if e.charactersDB != nil {
		e.charactersDB.Close()
	}
	if e.authDB != nil {
		e.authDB.Close()
	}
	if e.worldDB != nil {
		e.worldDB.Close()
	}
}

func (e *Exporter) collectPlayerMetrics() error {
	// Reset metrics
	playersOnline.Reset()
	playersTotal.Reset()
	playersByLevel.Reset()
	playersByClass.Reset()

	// Query for online players by faction
	query := `
		SELECT 
			race,
			COUNT(*) as count
		FROM characters 
		WHERE online = 1 
		GROUP BY race
	`
	rows, err := e.charactersDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var race int
		var count int
		if err := rows.Scan(&race, &count); err != nil {
			return err
		}
		faction := raceToFaction[race]
		if faction != "" {
			playersOnline.WithLabelValues(faction).Add(float64(count))
		}
	}

	// Query for total players by faction
	query = `
		SELECT 
			race,
			COUNT(*) as count
		FROM characters 
		GROUP BY race
	`
	rows, err = e.charactersDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var race int
		var count int
		if err := rows.Scan(&race, &count); err != nil {
			return err
		}
		faction := raceToFaction[race]
		if faction != "" {
			playersTotal.WithLabelValues(faction).Add(float64(count))
		}
	}

	// Query for players by level and faction
	query = `
		SELECT 
			level,
			race,
			COUNT(*) as count
		FROM characters 
		GROUP BY level, race
	`
	rows, err = e.charactersDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var level, race, count int
		if err := rows.Scan(&level, &race, &count); err != nil {
			return err
		}
		faction := raceToFaction[race]
		if faction != "" {
			playersByLevel.WithLabelValues(fmt.Sprintf("%d", level), faction).Add(float64(count))
		}
	}

	// Query for players by class and faction
	query = `
		SELECT 
			class,
			race,
			COUNT(*) as count
		FROM characters 
		GROUP BY class, race
	`
	rows, err = e.charactersDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var class, race, count int
		if err := rows.Scan(&class, &race, &count); err != nil {
			return err
		}
		faction := raceToFaction[race]
		className := classNames[class]
		if faction != "" && className != "" {
			playersByClass.WithLabelValues(className, faction).Add(float64(count))
		}
	}

	return nil
}

func (e *Exporter) collectMailMetrics() error {
	// Reset metrics
	mailTotal.Set(0)
	mailByFaction.Reset()
	mailWithItems.Set(0)

	// Query for total mail count
	var totalMail int
	query := `SELECT COUNT(*) FROM mail`
	err := e.charactersDB.QueryRow(query).Scan(&totalMail)
	if err != nil {
		return err
	}
	mailTotal.Add(float64(totalMail))

	// Query for mail with items
	var mailWithItemsCount int
	query = `SELECT COUNT(*) FROM mail WHERE has_items = 1`
	err = e.charactersDB.QueryRow(query).Scan(&mailWithItemsCount)
	if err != nil {
		return err
	}
	mailWithItems.Add(float64(mailWithItemsCount))

	// Query for mail by faction (based on sender's race)
	query = `
		SELECT 
			c.race,
			COUNT(*) as count
		FROM mail m
		JOIN characters c ON m.sender = c.guid
		GROUP BY c.race
	`
	rows, err := e.charactersDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var race, count int
		if err := rows.Scan(&race, &count); err != nil {
			return err
		}
		faction := raceToFaction[race]
		if faction != "" {
			mailByFaction.WithLabelValues(faction).Add(float64(count))
		}
	}

	return nil
}

func (e *Exporter) collectAccountMetrics() error {
	// Reset metrics
	accountsTotal.Set(0)
	accountsOnline.Set(0)
	accountsBanned.Set(0)

	// Query for total accounts (auth database)
	var totalAccounts int
	query := `SELECT COUNT(*) FROM account`
	err := e.authDB.QueryRow(query).Scan(&totalAccounts)
	if err != nil {
		return err
	}
	accountsTotal.Add(float64(totalAccounts))

	// Query for online accounts (auth database)
	var onlineAccounts int
	query = `SELECT COUNT(*) FROM account WHERE online = 1`
	err = e.authDB.QueryRow(query).Scan(&onlineAccounts)
	if err != nil {
		return err
	}
	accountsOnline.Add(float64(onlineAccounts))

	// Query for banned accounts (auth database)
	var bannedAccounts int
	query = `SELECT COUNT(*) FROM account_banned WHERE active = 1`
	err = e.authDB.QueryRow(query).Scan(&bannedAccounts)
	if err != nil {
		return err
	}
	accountsBanned.Add(float64(bannedAccounts))

	return nil
}

func (e *Exporter) collectServerMetrics() error {
	// Reset metrics
	serverUptime.Set(0)
	serverMaxPlayers.Set(0)

	// Query for server uptime and max players (auth database)
	query := `
		SELECT 
			uptime,
			maxplayers
		FROM uptime 
		WHERE realmid = 1 
		ORDER BY starttime DESC 
		LIMIT 1
	`
	var uptime, maxPlayers int
	err := e.authDB.QueryRow(query).Scan(&uptime, &maxPlayers)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err != sql.ErrNoRows {
		serverUptime.Add(float64(uptime))
		serverMaxPlayers.Add(float64(maxPlayers))
	}

	return nil
}

func (e *Exporter) collectAuctionMetrics() error {
	auctionCount.Reset()
	// houseid: 7 (neutral), 1 (alliance), 2 (horde)
	houseMap := map[int]string{1: "Alliance", 2: "Horde", 7: "Neutral"}
	query := `SELECT houseid, COUNT(*) FROM auctionhouse GROUP BY houseid`
	rows, err := e.charactersDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var houseid, count int
		if err := rows.Scan(&houseid, &count); err != nil {
			return err
		}
		house := houseMap[houseid]
		if house == "" {
			house = fmt.Sprintf("%d", houseid)
		}
		auctionCount.WithLabelValues(house).Set(float64(count))
	}
	return nil
}

func (e *Exporter) collectGuildMetrics() error {
	guildCount.Set(0)
	query := `SELECT COUNT(*) FROM guild`
	var count int
	err := e.charactersDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	guildCount.Set(float64(count))
	return nil
}

func (e *Exporter) collectMaxLevelCharMetrics() error {
	maxLevelCharCount.Reset()
	// AzerothCore WotLK max level is 80
	query := `SELECT race, COUNT(*) FROM characters WHERE level = 80 GROUP BY race`
	rows, err := e.charactersDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var race, count int
		if err := rows.Scan(&race, &count); err != nil {
			return err
		}
		faction := raceToFaction[race]
		if faction != "" {
			maxLevelCharCount.WithLabelValues(faction).Set(float64(count))
		}
	}
	return nil
}

func (e *Exporter) collectUnreadMailMetrics() error {
	unreadMailCount.Set(0)
	query := `SELECT COUNT(*) FROM mail WHERE checked = 0`
	var count int
	err := e.charactersDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	unreadMailCount.Set(float64(count))
	return nil
}

func (e *Exporter) collectGMAccountMetrics() error {
	gmAccountCount.Set(0)
	query := `SELECT COUNT(DISTINCT id) FROM account_access WHERE gmlevel > 0 AND RealmID IN (-1, 1)`
	var count int
	err := e.authDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	gmAccountCount.Set(float64(count))
	return nil
}

func (e *Exporter) collectLastServerRestartMetrics() error {
	lastServerRestart.Set(0)
	query := `SELECT starttime FROM uptime WHERE realmid = 1 ORDER BY starttime DESC LIMIT 1`
	var starttime int64
	err := e.authDB.QueryRow(query).Scan(&starttime)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil {
		lastServerRestart.Set(float64(starttime))
	}
	return nil
}

func (e *Exporter) collectBannedCharMetrics() error {
	bannedCharCount.Set(0)
	query := `SELECT COUNT(DISTINCT guid) FROM character_banned WHERE active = 1`
	var count int
	err := e.charactersDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	bannedCharCount.Set(float64(count))
	return nil
}

func (e *Exporter) collectChatMetrics() error {
	// Channel metrics (characters database)
	channelCount.Set(0)
	query := `SELECT COUNT(*) FROM channels`
	var count int
	err := e.charactersDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	channelCount.Set(float64(count))

	// Channel bans (characters database)
	channelBans.Set(0)
	query = `SELECT COUNT(*) FROM channels_bans`
	err = e.charactersDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	channelBans.Set(float64(count))

	// Log counts by type (auth database)
	logCountByType.Reset()
	query = `SELECT type, COUNT(*) FROM logs GROUP BY type`
	rows, err := e.authDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var logType string
		var count int
		if err := rows.Scan(&logType, &count); err != nil {
			return err
		}
		logCountByType.WithLabelValues(logType).Set(float64(count))
	}

	// Guild events (characters database)
	guildEventCount.Set(0)
	query = `SELECT COUNT(*) FROM guild_eventlog`
	err = e.charactersDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	guildEventCount.Set(float64(count))

	// Money logs (characters database)
	moneyLogCount.Set(0)
	query = `SELECT COUNT(*) FROM log_money`
	err = e.charactersDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	moneyLogCount.Set(float64(count))

	// Encounter logs (characters database)
	encounterLogCount.Set(0)
	query = `SELECT COUNT(*) FROM log_encounter`
	err = e.charactersDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	encounterLogCount.Set(float64(count))

	// Arena logs (characters database)
	arenaLogCount.Set(0)
	query = `SELECT COUNT(*) FROM log_arena_fights`
	err = e.charactersDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	arenaLogCount.Set(float64(count))

	// IP action logs (auth database)
	ipActionLogCount.Set(0)
	query = `SELECT COUNT(*) FROM logs_ip_actions`
	err = e.authDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	ipActionLogCount.Set(float64(count))

	return nil
}

func (e *Exporter) collectInstanceMetrics() error {
	// Active instances
	activeInstanceCount.Set(0)
	query := `SELECT COUNT(*) FROM instance WHERE resettime > UNIX_TIMESTAMP()`
	var count int
	err := e.charactersDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	activeInstanceCount.Set(float64(count))

	// Instances by difficulty
	instancesByDifficulty.Reset()
	query = `SELECT difficulty, COUNT(*) FROM instance GROUP BY difficulty`
	rows, err := e.charactersDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var difficulty, count int
		if err := rows.Scan(&difficulty, &count); err != nil {
			return err
		}
		difficultyName := getDifficultyName(difficulty)
		instancesByDifficulty.WithLabelValues(difficultyName).Set(float64(count))
	}

	// Completed encounters
	completedEncounters.Reset()
	query = `SELECT id, completedEncounters FROM instance WHERE completedEncounters > 0`
	rows, err = e.charactersDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var instanceID, encounters int
		if err := rows.Scan(&instanceID, &encounters); err != nil {
			return err
		}
		completedEncounters.WithLabelValues(fmt.Sprintf("%d", instanceID)).Set(float64(encounters))
	}

	// Instance resets
	instanceResets.Reset()
	query = `SELECT mapid, difficulty, resettime FROM instance_reset`
	rows, err = e.charactersDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var mapID, difficulty, resetTime int
		if err := rows.Scan(&mapID, &difficulty, &resetTime); err != nil {
			return err
		}
		difficultyName := getDifficultyName(difficulty)
		instanceResets.WithLabelValues(fmt.Sprintf("%d", mapID), difficultyName).Set(float64(resetTime))
	}

	// Characters in instances
	charactersInInstances.Set(0)
	query = `SELECT COUNT(DISTINCT guid) FROM character_instance`
	err = e.charactersDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	charactersInInstances.Set(float64(count))

	// LFG data
	lfgDataCount.Reset()
	query = `SELECT state, COUNT(*) FROM lfg_data GROUP BY state`
	rows, err = e.charactersDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var state, count int
		if err := rows.Scan(&state, &count); err != nil {
			return err
		}
		stateName := getLFGStateName(state)
		lfgDataCount.WithLabelValues(stateName).Set(float64(count))
	}

	// Lag reports
	lagReportsCount.Set(0)
	query = `SELECT COUNT(*) FROM lag_reports`
	err = e.charactersDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	lagReportsCount.Set(float64(count))

	// Instance saves
	instanceSavesCount.Set(0)
	query = `SELECT COUNT(*) FROM instance_saved_go_state_data`
	err = e.charactersDB.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	instanceSavesCount.Set(float64(count))

	return nil
}

// Helper functions for readable names
func getDifficultyName(difficulty int) string {
	difficultyNames := map[int]string{
		0: "Normal",
		1: "Heroic",
		2: "10_Player",
		3: "25_Player",
		4: "10_Player_Heroic",
		5: "25_Player_Heroic",
	}
	if name, exists := difficultyNames[difficulty]; exists {
		return name
	}
	return fmt.Sprintf("Unknown_%d", difficulty)
}

func getLFGStateName(state int) string {
	stateNames := map[int]string{
		0: "None",
		1: "RoleCheck",
		2: "Queued",
		3: "Proposal",
		4: "Boot",
		5: "Dungeon",
		6: "FinishedDungeon",
	}
	if name, exists := stateNames[state]; exists {
		return name
	}
	return fmt.Sprintf("Unknown_%d", state)
}

func (e *Exporter) collectNetworkMetrics() error {
	// Player latency statistics
	playerLatencyStats.Reset()

	// Average latency
	var avgLatency float64
	query := `SELECT AVG(latency) FROM characters WHERE online = 1 AND latency > 0`
	err := e.charactersDB.QueryRow(query).Scan(&avgLatency)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err != sql.ErrNoRows {
		playerLatencyStats.WithLabelValues("average").Set(avgLatency)
		averageLatency.Set(avgLatency)
	}

	// High latency players (>200ms)
	var highLatencyCount int
	query = `SELECT COUNT(*) FROM characters WHERE online = 1 AND latency > 200`
	err = e.charactersDB.QueryRow(query).Scan(&highLatencyCount)
	if err != nil {
		return err
	}
	highLatencyPlayers.Set(float64(highLatencyCount))
	playerLatencyStats.WithLabelValues("high_latency").Set(float64(highLatencyCount))

	// Min/Max latency
	var minLatency, maxLatency int
	query = `SELECT MIN(latency), MAX(latency) FROM characters WHERE online = 1 AND latency > 0`
	err = e.charactersDB.QueryRow(query).Scan(&minLatency, &maxLatency)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err != sql.ErrNoRows {
		playerLatencyStats.WithLabelValues("min").Set(float64(minLatency))
		playerLatencyStats.WithLabelValues("max").Set(float64(maxLatency))
	}

	// IP bans
	ipBannedCount.Set(0)
	var bannedCount int
	query = `SELECT COUNT(*) FROM ip_banned`
	err = e.authDB.QueryRow(query).Scan(&bannedCount)
	if err != nil {
		return err
	}
	ipBannedCount.Set(float64(bannedCount))

	// IP action logs by type
	ipActionLogsByType.Reset()
	query = `SELECT type, COUNT(*) FROM logs_ip_actions GROUP BY type`
	rows, err := e.authDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var actionType, count int
		if err := rows.Scan(&actionType, &count); err != nil {
			return err
		}
		typeName := getIPActionTypeName(actionType)
		ipActionLogsByType.WithLabelValues(typeName).Set(float64(count))
	}

	// Lag reports by type
	lagReportsByType.Reset()
	query = `SELECT lagType, COUNT(*) FROM lag_reports GROUP BY lagType`
	rows, err = e.charactersDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var lagType, count int
		if err := rows.Scan(&lagType, &count); err != nil {
			return err
		}
		lagTypeName := getLagTypeName(lagType)
		lagReportsByType.WithLabelValues(lagTypeName).Set(float64(count))
	}

	// Network activity by IP (top 10 most active IPs)
	networkActivityByIP.Reset()
	query = `SELECT ip, COUNT(*) as activity FROM logs_ip_actions GROUP BY ip ORDER BY activity DESC LIMIT 10`
	rows, err = e.authDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var ip string
		var activity int
		if err := rows.Scan(&ip, &activity); err != nil {
			return err
		}
		networkActivityByIP.WithLabelValues(ip).Set(float64(activity))
	}

	return nil
}

// Helper functions for networking
func getIPActionTypeName(actionType int) string {
	typeNames := map[int]string{
		0: "Login",
		1: "Failed_Login",
		2: "Logout",
		3: "Character_Create",
		4: "Character_Delete",
		5: "Character_Login",
		6: "Character_Logout",
		7: "Password_Change",
		8: "Account_Create",
		9: "Account_Delete",
	}
	if name, exists := typeNames[actionType]; exists {
		return name
	}
	return fmt.Sprintf("Unknown_%d", actionType)
}

func getLagTypeName(lagType int) string {
	lagTypeNames := map[int]string{
		0: "World",
		1: "Instance",
		2: "Battleground",
		3: "Arena",
		4: "Raid",
	}
	if name, exists := lagTypeNames[lagType]; exists {
		return name
	}
	return fmt.Sprintf("Unknown_%d", lagType)
}

func getBattlegroundTypeName(bgType int) string {
	bgTypeNames := map[int]string{
		1:  "Alterac Valley",
		2:  "Warsong Gulch",
		3:  "Arathi Basin",
		4:  "Eye of the Storm",
		5:  "Strand of the Ancients",
		6:  "Isle of Conquest",
		7:  "Twin Peaks",
		8:  "Battle for Gilneas",
		9:  "Temple of Kotmogu",
		10: "Silvershard Mines",
		11: "Deepwind Gorge",
	}
	if name, exists := bgTypeNames[bgType]; exists {
		return name
	}
	return fmt.Sprintf("Unknown_%d", bgType)
}

func getDesertionTypeName(desertionType int) string {
	desertionTypeNames := map[int]string{
		0: "Leave",
		1: "Offline",
		2: "Desert",
		3: "Finish",
	}
	if name, exists := desertionTypeNames[desertionType]; exists {
		return name
	}
	return fmt.Sprintf("Unknown_%d", desertionType)
}

func (e *Exporter) collectBattlegroundMetrics() error {
	// Reset vector metrics
	battlegroundDesertersByType.Reset()
	battlegroundStats.Reset()
	battlegroundsByType.Reset()
	battlegroundsByBracket.Reset()
	battlegroundWinsByFaction.Reset()
	battlegroundPlayerStats.Reset()
	battlegroundTemplates.Reset()
	battlegroundTemplateDetails.Reset()
	recentBattlegrounds.Reset()

	// Battleground deserters
	var deserterCount int
	err := e.charactersDB.QueryRow("SELECT COUNT(*) FROM battleground_deserters").Scan(&deserterCount)
	if err != nil {
		return fmt.Errorf("error querying battleground deserters: %v", err)
	}
	battlegroundDeserters.Set(float64(deserterCount))

	// Battleground deserters by type
	rows, err := e.charactersDB.Query(`
		SELECT type, COUNT(*) as count 
		FROM battleground_deserters 
		GROUP BY type
	`)
	if err != nil {
		return fmt.Errorf("error querying battleground deserters by type: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var desertionType int
		var count int
		if err := rows.Scan(&desertionType, &count); err != nil {
			continue
		}
		battlegroundDesertersByType.WithLabelValues(getDesertionTypeName(desertionType)).Set(float64(count))
	}

	// Random battleground queue
	var queueCount int
	err = e.charactersDB.QueryRow("SELECT COUNT(*) FROM character_battleground_random").Scan(&queueCount)
	if err != nil {
		return fmt.Errorf("error querying random battleground queue: %v", err)
	}
	randomBattlegroundQueue.Set(float64(queueCount))

	// Battleground statistics
	var totalBattlegrounds, totalPlayers int
	err = e.charactersDB.QueryRow(`
		SELECT 
			COUNT(DISTINCT id) as total_battlegrounds,
			COUNT(DISTINCT character_guid) as total_players
		FROM pvpstats_battlegrounds bg
		LEFT JOIN pvpstats_players bp ON bg.id = bp.battleground_id
	`).Scan(&totalBattlegrounds, &totalPlayers)
	if err != nil {
		return fmt.Errorf("error querying battleground stats: %v", err)
	}
	battlegroundStats.WithLabelValues("total_battlegrounds").Set(float64(totalBattlegrounds))
	battlegroundStats.WithLabelValues("total_players").Set(float64(totalPlayers))

	// Battlegrounds by type
	rows, err = e.charactersDB.Query(`
		SELECT type, COUNT(*) as count 
		FROM pvpstats_battlegrounds 
		GROUP BY type
	`)
	if err != nil {
		return fmt.Errorf("error querying battlegrounds by type: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var bgType int
		var count int
		if err := rows.Scan(&bgType, &count); err != nil {
			continue
		}
		battlegroundsByType.WithLabelValues(getBattlegroundTypeName(bgType)).Set(float64(count))
	}

	// Battlegrounds by bracket
	rows, err = e.charactersDB.Query(`
		SELECT bracket_id, COUNT(*) as count 
		FROM pvpstats_battlegrounds 
		GROUP BY bracket_id
	`)
	if err != nil {
		return fmt.Errorf("error querying battlegrounds by bracket: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var bracket int
		var count int
		if err := rows.Scan(&bracket, &count); err != nil {
			continue
		}
		battlegroundsByBracket.WithLabelValues(fmt.Sprintf("bracket_%d", bracket)).Set(float64(count))
	}

	// Battleground wins by faction
	rows, err = e.charactersDB.Query(`
		SELECT winner_faction, COUNT(*) as count 
		FROM pvpstats_battlegrounds 
		WHERE winner_faction IN (0, 1)
		GROUP BY winner_faction
	`)
	if err != nil {
		return fmt.Errorf("error querying battleground wins by faction: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var faction int
		var count int
		if err := rows.Scan(&faction, &count); err != nil {
			continue
		}
		factionName := "Horde"
		if faction == 0 {
			factionName = "Alliance"
		}
		battlegroundWinsByFaction.WithLabelValues(factionName).Set(float64(count))
	}

	// Battleground player statistics
	rows, err = e.charactersDB.Query(`
		SELECT 
			COUNT(*) as total_participants,
			SUM(CASE WHEN winner = 1 THEN 1 ELSE 0 END) as total_winners,
			AVG(score_killing_blows) as avg_killing_blows,
			AVG(score_deaths) as avg_deaths,
			AVG(score_honorable_kills) as avg_honorable_kills,
			AVG(score_bonus_honor) as avg_bonus_honor,
			AVG(score_damage_done) as avg_damage_done,
			AVG(score_healing_done) as avg_healing_done
		FROM pvpstats_players
	`)
	if err != nil {
		return fmt.Errorf("error querying battleground player stats: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		var totalParticipants, totalWinners int
		var avgKillingBlows, avgDeaths, avgHonorableKills, avgBonusHonor, avgDamageDone, avgHealingDone float64
		if err := rows.Scan(&totalParticipants, &totalWinners, &avgKillingBlows, &avgDeaths, &avgHonorableKills, &avgBonusHonor, &avgDamageDone, &avgHealingDone); err == nil {
			battlegroundPlayerStats.WithLabelValues("total_participants").Set(float64(totalParticipants))
			battlegroundPlayerStats.WithLabelValues("total_winners").Set(float64(totalWinners))
			battlegroundPlayerStats.WithLabelValues("avg_killing_blows").Set(avgKillingBlows)
			battlegroundPlayerStats.WithLabelValues("avg_deaths").Set(avgDeaths)
			battlegroundPlayerStats.WithLabelValues("avg_honorable_kills").Set(avgHonorableKills)
			battlegroundPlayerStats.WithLabelValues("avg_bonus_honor").Set(avgBonusHonor)
			battlegroundPlayerStats.WithLabelValues("avg_damage_done").Set(avgDamageDone)
			battlegroundPlayerStats.WithLabelValues("avg_healing_done").Set(avgHealingDone)
		}
	}

	// Battleground templates
	rows, err = e.worldDB.Query(`
		SELECT ID, ScriptName, Comment, MinPlayersPerTeam, MaxPlayersPerTeam, MinLvl, MaxLvl, Weight
		FROM battleground_template
	`)
	if err != nil {
		return fmt.Errorf("error querying battleground templates: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, minPlayers, maxPlayers, minLvl, maxLvl, weight int
		var scriptName, comment string
		if err := rows.Scan(&id, &scriptName, &comment, &minPlayers, &maxPlayers, &minLvl, &maxLvl, &weight); err != nil {
			continue
		}
		
		// Create a more descriptive label
		label := scriptName
		if label == "" {
			label = comment
		}
		if label == "" {
			label = fmt.Sprintf("BG_%d", id)
		}
		
		battlegroundTemplates.WithLabelValues(fmt.Sprintf("%d", id), label).Set(float64(weight))
		
		// Add detailed template information
		battlegroundTemplateDetails.WithLabelValues(
			fmt.Sprintf("%d", id),
			label,
			fmt.Sprintf("%d", minLvl),
			fmt.Sprintf("%d", maxLvl),
			fmt.Sprintf("%d", minPlayers),
			fmt.Sprintf("%d", maxPlayers),
		).Set(float64(weight))
	}

	// Recent battleground activity (last 24 hours, 7 days, 30 days)
	rows, err = e.charactersDB.Query(`
		SELECT 
			SUM(CASE WHEN date >= DATE_SUB(NOW(), INTERVAL 24 HOUR) THEN 1 ELSE 0 END) as last_24h,
			SUM(CASE WHEN date >= DATE_SUB(NOW(), INTERVAL 7 DAY) THEN 1 ELSE 0 END) as last_7d,
			SUM(CASE WHEN date >= DATE_SUB(NOW(), INTERVAL 30 DAY) THEN 1 ELSE 0 END) as last_30d
		FROM pvpstats_battlegrounds
	`)
	if err != nil {
		return fmt.Errorf("error querying recent battleground activity: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		var last24h, last7d, last30d int
		if err := rows.Scan(&last24h, &last7d, &last30d); err == nil {
			recentBattlegrounds.WithLabelValues("last_24h").Set(float64(last24h))
			recentBattlegrounds.WithLabelValues("last_7d").Set(float64(last7d))
			recentBattlegrounds.WithLabelValues("last_30d").Set(float64(last30d))
		}
	}

	return nil
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Collect all metrics
	if err := e.collectPlayerMetrics(); err != nil {
		log.Printf("Error collecting player metrics: %v", err)
	}
	if err := e.collectMailMetrics(); err != nil {
		log.Printf("Error collecting mail metrics: %v", err)
	}
	if err := e.collectAccountMetrics(); err != nil {
		log.Printf("Error collecting account metrics: %v", err)
	}
	if err := e.collectServerMetrics(); err != nil {
		log.Printf("Error collecting server metrics: %v", err)
	}
	if err := e.collectAuctionMetrics(); err != nil {
		log.Printf("Error collecting auction metrics: %v", err)
	}
	if err := e.collectGuildMetrics(); err != nil {
		log.Printf("Error collecting guild metrics: %v", err)
	}
	if err := e.collectMaxLevelCharMetrics(); err != nil {
		log.Printf("Error collecting max level char metrics: %v", err)
	}
	if err := e.collectUnreadMailMetrics(); err != nil {
		log.Printf("Error collecting unread mail metrics: %v", err)
	}
	if err := e.collectGMAccountMetrics(); err != nil {
		log.Printf("Error collecting GM account metrics: %v", err)
	}
	if err := e.collectLastServerRestartMetrics(); err != nil {
		log.Printf("Error collecting last server restart metrics: %v", err)
	}
	if err := e.collectBannedCharMetrics(); err != nil {
		log.Printf("Error collecting banned char metrics: %v", err)
	}
	if err := e.collectChatMetrics(); err != nil {
		log.Printf("Error collecting chat metrics: %v", err)
	}
	if err := e.collectInstanceMetrics(); err != nil {
		log.Printf("Error collecting instance metrics: %v", err)
	}
	if err := e.collectNetworkMetrics(); err != nil {
		log.Printf("Error collecting network metrics: %v", err)
	}
	if err := e.collectBattlegroundMetrics(); err != nil {
		log.Printf("Error collecting battleground metrics: %v", err)
	}

	// Send all metrics
	playersOnline.Collect(ch)
	playersTotal.Collect(ch)
	playersByLevel.Collect(ch)
	playersByClass.Collect(ch)
	mailTotal.Collect(ch)
	mailByFaction.Collect(ch)
	mailWithItems.Collect(ch)
	accountsTotal.Collect(ch)
	accountsOnline.Collect(ch)
	accountsBanned.Collect(ch)
	serverUptime.Collect(ch)
	serverMaxPlayers.Collect(ch)
	auctionCount.Collect(ch)
	guildCount.Collect(ch)
	maxLevelCharCount.Collect(ch)
	unreadMailCount.Collect(ch)
	gmAccountCount.Collect(ch)
	lastServerRestart.Collect(ch)
	bannedCharCount.Collect(ch)
	channelCount.Collect(ch)
	channelBans.Collect(ch)
	logCountByType.Collect(ch)
	guildEventCount.Collect(ch)
	moneyLogCount.Collect(ch)
	encounterLogCount.Collect(ch)
	arenaLogCount.Collect(ch)
	ipActionLogCount.Collect(ch)
	activeInstanceCount.Collect(ch)
	instancesByDifficulty.Collect(ch)
	completedEncounters.Collect(ch)
	instanceResets.Collect(ch)
	charactersInInstances.Collect(ch)
	lfgDataCount.Collect(ch)
	lagReportsCount.Collect(ch)
	instanceSavesCount.Collect(ch)
	playerLatencyStats.Collect(ch)
	ipBannedCount.Collect(ch)
	ipActionLogsByType.Collect(ch)
	lagReportsByType.Collect(ch)
	averageLatency.Collect(ch)
	highLatencyPlayers.Collect(ch)
	networkActivityByIP.Collect(ch)
	battlegroundDeserters.Collect(ch)
	battlegroundDesertersByType.Collect(ch)
	randomBattlegroundQueue.Collect(ch)
	battlegroundStats.Collect(ch)
	battlegroundsByType.Collect(ch)
	battlegroundsByBracket.Collect(ch)
	battlegroundWinsByFaction.Collect(ch)
	battlegroundPlayerStats.Collect(ch)
	battlegroundTemplates.Collect(ch)
	battlegroundTemplateDetails.Collect(ch)
	recentBattlegrounds.Collect(ch)
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	playersOnline.Describe(ch)
	playersTotal.Describe(ch)
	playersByLevel.Describe(ch)
	playersByClass.Describe(ch)
	mailTotal.Describe(ch)
	mailByFaction.Describe(ch)
	mailWithItems.Describe(ch)
	accountsTotal.Describe(ch)
	accountsOnline.Describe(ch)
	accountsBanned.Describe(ch)
	serverUptime.Describe(ch)
	serverMaxPlayers.Describe(ch)
	auctionCount.Describe(ch)
	guildCount.Describe(ch)
	maxLevelCharCount.Describe(ch)
	unreadMailCount.Describe(ch)
	gmAccountCount.Describe(ch)
	lastServerRestart.Describe(ch)
	bannedCharCount.Describe(ch)
	channelCount.Describe(ch)
	channelBans.Describe(ch)
	logCountByType.Describe(ch)
	guildEventCount.Describe(ch)
	moneyLogCount.Describe(ch)
	encounterLogCount.Describe(ch)
	arenaLogCount.Describe(ch)
	ipActionLogCount.Describe(ch)
	activeInstanceCount.Describe(ch)
	instancesByDifficulty.Describe(ch)
	completedEncounters.Describe(ch)
	instanceResets.Describe(ch)
	charactersInInstances.Describe(ch)
	lfgDataCount.Describe(ch)
	lagReportsCount.Describe(ch)
	instanceSavesCount.Describe(ch)
	playerLatencyStats.Describe(ch)
	ipBannedCount.Describe(ch)
	ipActionLogsByType.Describe(ch)
	lagReportsByType.Describe(ch)
	averageLatency.Describe(ch)
	highLatencyPlayers.Describe(ch)
	networkActivityByIP.Describe(ch)
	battlegroundDeserters.Describe(ch)
	battlegroundDesertersByType.Describe(ch)
	randomBattlegroundQueue.Describe(ch)
	battlegroundStats.Describe(ch)
	battlegroundsByType.Describe(ch)
	battlegroundsByBracket.Describe(ch)
	battlegroundWinsByFaction.Describe(ch)
	battlegroundPlayerStats.Describe(ch)
	battlegroundTemplates.Describe(ch)
	battlegroundTemplateDetails.Describe(ch)
	recentBattlegrounds.Describe(ch)
}

// Helper to build DSN from env vars
func buildDSN() string {
	user := os.Getenv("WOW_DB_USER")
	pass := os.Getenv("WOW_DB_PASS")
	host := os.Getenv("WOW_DB_HOST")
	port := os.Getenv("WOW_DB_PORT")
	if port == "" {
		port = "3306"
	}
	// Default to acore_characters DB for connection
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/acore_characters?parseTime=true", user, pass, host, port)
}

func main() {
	// Get database connection string from environment
	dsn := os.Getenv("WOW_DB_DSN")
	if dsn == "" {
		dsn = buildDSN()
	}

	// Create exporter
	exporter, err := NewExporter(dsn)
	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}
	defer exporter.Close()

	// Register metrics
	prometheus.MustRegister(exporter)

	// Get port from environment, default to 7000
	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}

	// Setup HTTP server
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
			<head><title>WoW Private Server Exporter</title></head>
			<body>
				<h1>WoW Private Server Exporter</h1>
				<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>
		`))
	})

	log.Printf("Starting WoW Private Server Exporter on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
