package exporter

import (
	"database/sql"
	"fmt"

	"github.com/scottjab/prom-azerothcore-exporter/metrics"
	"github.com/scottjab/prom-azerothcore-exporter/pkg/constants"
	"github.com/scottjab/prom-azerothcore-exporter/pkg/database"
)

// Collection methods for the Exporter

func (e *Exporter) collectPlayerMetrics() error {
	// Reset metrics
	metrics.PlayersOnline.Reset()
	metrics.PlayersTotal.Reset()
	metrics.PlayersByLevel.Reset()
	metrics.PlayersByClass.Reset()
	metrics.OnlinePlayersByLevel.Reset()

	// Query for online players by faction
	query := `
		SELECT 
			race,
			COUNT(*) as count
		FROM characters 
		WHERE online = 1 
		GROUP BY race
	`
	rows, err := e.connections.Characters.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)

	for rows.Next() {
		var race int
		var count int
		if err := rows.Scan(&race, &count); err != nil {
			return err
		}
		faction := constants.RaceToFaction[race]
		if faction != "" {
			metrics.PlayersOnline.WithLabelValues(faction).Add(float64(count))
		}
	}

	// Query for total players by faction
	// Exclude likely test characters: very recent creations and test names
	query = `
		SELECT 
			race,
			COUNT(*) as count
		FROM characters 
		WHERE (deleteDate IS NULL OR deleteDate = 0)
		AND (creation_date IS NULL OR creation_date < DATE_SUB(NOW(), INTERVAL 24 HOUR)) -- Exclude characters created in last 24 hours
		AND name NOT LIKE '%test%'
		AND name NOT LIKE '%admin%'
		AND name NOT LIKE '%gm%'
		AND name NOT LIKE '%dev%'
		GROUP BY race
	`
	rows, err = e.connections.Characters.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)

	for rows.Next() {
		var race int
		var count int
		if err := rows.Scan(&race, &count); err != nil {
			return err
		}
		faction := constants.RaceToFaction[race]
		if faction != "" {
			metrics.PlayersTotal.WithLabelValues(faction).Add(float64(count))
		}
	}

	// Query for players by level and faction
	// Exclude likely test characters: filter by creation date, test names, and never logged in characters
	query = `
		SELECT 
			level,
			race,
			COUNT(*) as count
		FROM characters 
		WHERE (deleteDate IS NULL OR deleteDate = 0)
		AND (creation_date IS NULL OR creation_date < DATE_SUB(NOW(), INTERVAL 24 HOUR)) -- Exclude characters created in last 24 hours
		AND logout_time > 0 -- Exclude characters that have never logged in
		AND name NOT LIKE '%test%'
		AND name NOT LIKE '%admin%'
		AND name NOT LIKE '%gm%'
		AND name NOT LIKE '%dev%'
		AND name NOT LIKE '%temp%'
		AND name NOT LIKE '%demo%'
		AND name NOT LIKE '%example%'
		GROUP BY level, race
	`
	rows, err = e.connections.Characters.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)

	for rows.Next() {
		var level, race, count int
		if err := rows.Scan(&level, &race, &count); err != nil {
			return err
		}
		faction := constants.RaceToFaction[race]
		if faction != "" {
			metrics.PlayersByLevel.WithLabelValues(fmt.Sprintf("%d", level), faction).Add(float64(count))
		}
	}

	// Query for players by class and faction
	// Exclude likely test characters: very recent creations and test names
	query = `
		SELECT 
			class,
			race,
			COUNT(*) as count
		FROM characters 
		WHERE (deleteDate IS NULL OR deleteDate = 0)
		AND (creation_date IS NULL OR creation_date < DATE_SUB(NOW(), INTERVAL 24 HOUR)) -- Exclude characters created in last 24 hours
		AND name NOT LIKE '%test%'
		AND name NOT LIKE '%admin%'
		AND name NOT LIKE '%gm%'
		AND name NOT LIKE '%dev%'
		GROUP BY class, race
	`
	rows, err = e.connections.Characters.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)

	for rows.Next() {
		var class, race, count int
		if err := rows.Scan(&class, &race, &count); err != nil {
			return err
		}
		faction := constants.RaceToFaction[race]
		className := constants.ClassNames[class]
		if faction != "" && className != "" {
			metrics.PlayersByClass.WithLabelValues(className, faction).Add(float64(count))
		}
	}

	// Query for online players by level with character and account name
	// First get character information from characters database
	query = `
		SELECT 
			c.name,
			c.level,
			c.account
		FROM characters c
		WHERE c.online = 1
		AND (c.deleteDate IS NULL OR c.deleteDate = 0)
		ORDER BY c.level, c.name
	`
	rows, err = e.connections.Characters.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)

	// Create a map to store account IDs and their usernames
	accountMap := make(map[int]string)

	for rows.Next() {
		var characterName string
		var level, accountID int
		if err := rows.Scan(&characterName, &level, &accountID); err != nil {
			return err
		}

		// Get account username if we haven't already
		accountName, exists := accountMap[accountID]
		if !exists {
			// Query auth database for account username
			var username string
			authQuery := `SELECT username FROM account WHERE id = ?`
			err := e.connections.Auth.QueryRow(authQuery, accountID).Scan(&username)
			if err != nil {
				// If we can't get the username, use a placeholder
				username = fmt.Sprintf("account_%d", accountID)
			}
			accountMap[accountID] = username
			accountName = username
		}

		metrics.OnlinePlayersByLevel.WithLabelValues(
			characterName,
			accountName,
		).Set(float64(level))
	}

	return nil
}

func (e *Exporter) collectMailMetrics() error {
	// Reset metrics
	metrics.MailTotal.Set(0)
	metrics.MailByFaction.Reset()
	metrics.MailWithItems.Set(0)

	// Query for total mail count
	var totalMail int
	query := `SELECT COUNT(*) FROM mail`
	err := e.connections.Characters.QueryRow(query).Scan(&totalMail)
	if err != nil {
		return err
	}
	metrics.MailTotal.Add(float64(totalMail))

	// Query for mail with items
	var mailWithItemsCount int
	query = `SELECT COUNT(*) FROM mail WHERE has_items = 1`
	err = e.connections.Characters.QueryRow(query).Scan(&mailWithItemsCount)
	if err != nil {
		return err
	}
	metrics.MailWithItems.Add(float64(mailWithItemsCount))

	// Query for mail by faction (based on sender's race)
	// Exclude likely test characters: very recent creations and test names
	query = `
		SELECT 
			c.race,
			COUNT(*) as count
		FROM mail m
		JOIN characters c ON m.sender = c.guid
		WHERE (c.deleteDate IS NULL OR c.deleteDate = 0)
		AND (c.creation_date IS NULL OR c.creation_date < DATE_SUB(NOW(), INTERVAL 24 HOUR)) -- Exclude characters created in last 24 hours
		AND c.name NOT LIKE '%test%'
		AND c.name NOT LIKE '%admin%'
		AND c.name NOT LIKE '%gm%'
		AND c.name NOT LIKE '%dev%'
		GROUP BY c.race
	`
	rows, err := e.connections.Characters.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)

	for rows.Next() {
		var race, count int
		if err := rows.Scan(&race, &count); err != nil {
			return err
		}
		faction := constants.RaceToFaction[race]
		if faction != "" {
			metrics.MailByFaction.WithLabelValues(faction).Add(float64(count))
		}
	}

	return nil
}

func (e *Exporter) collectAccountMetrics() error {
	// Reset metrics
	metrics.AccountsTotal.Set(0)
	metrics.AccountsOnline.Set(0)
	metrics.AccountsBanned.Set(0)

	// Query for total accounts (auth database)
	var totalAccounts int
	query := `SELECT COUNT(*) FROM account`
	err := e.connections.Auth.QueryRow(query).Scan(&totalAccounts)
	if err != nil {
		return err
	}
	metrics.AccountsTotal.Add(float64(totalAccounts))

	// Query for online accounts (auth database)
	var onlineAccounts int
	query = `SELECT COUNT(*) FROM account WHERE online = 1`
	err = e.connections.Auth.QueryRow(query).Scan(&onlineAccounts)
	if err != nil {
		return err
	}
	metrics.AccountsOnline.Add(float64(onlineAccounts))

	// Query for banned accounts (auth database)
	var bannedAccounts int
	query = `SELECT COUNT(*) FROM account_banned WHERE active = 1`
	err = e.connections.Auth.QueryRow(query).Scan(&bannedAccounts)
	if err != nil {
		return err
	}
	metrics.AccountsBanned.Add(float64(bannedAccounts))

	return nil
}

func (e *Exporter) collectServerMetrics() error {
	// Reset metrics
	metrics.ServerUptime.Set(0)
	metrics.ServerMaxPlayers.Set(0)

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
	err := e.connections.Auth.QueryRow(query).Scan(&uptime, &maxPlayers)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err != sql.ErrNoRows {
		metrics.ServerUptime.Add(float64(uptime))
		metrics.ServerMaxPlayers.Add(float64(maxPlayers))
	}

	return nil
}

func (e *Exporter) collectAuctionMetrics() error {
	metrics.AuctionCount.Reset()
	// houseid: 7 (neutral), 1 (alliance), 2 (horde)
	houseMap := map[int]string{1: "Alliance", 2: "Horde", 7: "Neutral"}
	query := `SELECT houseid, COUNT(*) FROM auctionhouse GROUP BY houseid`
	rows, err := e.connections.Characters.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)
	for rows.Next() {
		var houseid, count int
		if err := rows.Scan(&houseid, &count); err != nil {
			return err
		}
		house := houseMap[houseid]
		if house == "" {
			house = fmt.Sprintf("%d", houseid)
		}
		metrics.AuctionCount.WithLabelValues(house).Set(float64(count))
	}
	return nil
}

func (e *Exporter) collectGuildMetrics() error {
	metrics.GuildCount.Set(0)
	query := `SELECT COUNT(*) FROM guild`
	var count int
	err := e.connections.Characters.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	metrics.GuildCount.Set(float64(count))
	return nil
}

func (e *Exporter) collectMaxLevelCharMetrics() error {
	metrics.MaxLevelCharCount.Reset()
	// AzerothCore WotLK max level is 80
	// Exclude likely test characters: very recent creations and test names
	query := `
		SELECT race, COUNT(*) 
		FROM characters 
		WHERE level = 80 
		AND (deleteDate IS NULL OR deleteDate = 0)
		AND (creation_date IS NULL OR creation_date < DATE_SUB(NOW(), INTERVAL 24 HOUR)) -- Exclude characters created in last 24 hours
		AND name NOT LIKE '%test%'
		AND name NOT LIKE '%admin%'
		AND name NOT LIKE '%gm%'
		AND name NOT LIKE '%dev%'
		GROUP BY race
	`
	rows, err := e.connections.Characters.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)
	for rows.Next() {
		var race, count int
		if err := rows.Scan(&race, &count); err != nil {
			return err
		}
		faction := constants.RaceToFaction[race]
		if faction != "" {
			metrics.MaxLevelCharCount.WithLabelValues(faction).Set(float64(count))
		}
	}
	return nil
}

func (e *Exporter) collectUnreadMailMetrics() error {
	metrics.UnreadMailCount.Set(0)
	query := `SELECT COUNT(*) FROM mail WHERE checked = 0`
	var count int
	err := e.connections.Characters.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	metrics.UnreadMailCount.Set(float64(count))
	return nil
}

func (e *Exporter) collectGMAccountMetrics() error {
	metrics.GMAccountCount.Set(0)
	query := `SELECT COUNT(DISTINCT id) FROM account_access WHERE gmlevel > 0 AND RealmID IN (-1, 1)`
	var count int
	err := e.connections.Auth.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	metrics.GMAccountCount.Set(float64(count))
	return nil
}

func (e *Exporter) collectLastServerRestartMetrics() error {
	metrics.LastServerRestart.Set(0)
	query := `SELECT starttime FROM uptime WHERE realmid = 1 ORDER BY starttime DESC LIMIT 1`
	var starttime int64
	err := e.connections.Auth.QueryRow(query).Scan(&starttime)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil {
		metrics.LastServerRestart.Set(float64(starttime))
	}
	return nil
}

func (e *Exporter) collectBannedCharMetrics() error {
	metrics.BannedCharCount.Set(0)
	query := `SELECT COUNT(DISTINCT guid) FROM character_banned WHERE active = 1`
	var count int
	err := e.connections.Characters.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	metrics.BannedCharCount.Set(float64(count))
	return nil
}

func (e *Exporter) collectChatMetrics() error {
	// Channel metrics (characters database)
	metrics.ChannelCount.Set(0)
	query := `SELECT COUNT(*) FROM channels`
	var count int
	err := e.connections.Characters.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	metrics.ChannelCount.Set(float64(count))

	// Channel bans (characters database)
	metrics.ChannelBans.Set(0)
	query = `SELECT COUNT(*) FROM channels_bans`
	err = e.connections.Characters.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	metrics.ChannelBans.Set(float64(count))

	// Log counts by type (auth database)
	metrics.LogCountByType.Reset()
	query = `SELECT type, COUNT(*) FROM logs GROUP BY type`
	rows, err := e.connections.Auth.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)
	for rows.Next() {
		var logType string
		var count int
		if err := rows.Scan(&logType, &count); err != nil {
			return err
		}
		metrics.LogCountByType.WithLabelValues(logType).Set(float64(count))
	}

	// Guild events (characters database)
	metrics.GuildEventCount.Set(0)
	query = `SELECT COUNT(*) FROM guild_eventlog`
	err = e.connections.Characters.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	metrics.GuildEventCount.Set(float64(count))

	// Money logs (characters database)
	metrics.MoneyLogCount.Set(0)
	query = `SELECT COUNT(*) FROM log_money`
	err = e.connections.Characters.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	metrics.MoneyLogCount.Set(float64(count))

	// Encounter logs (characters database)
	metrics.EncounterLogCount.Set(0)
	query = `SELECT COUNT(*) FROM log_encounter`
	err = e.connections.Characters.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	metrics.EncounterLogCount.Set(float64(count))

	// Arena logs (characters database)
	metrics.ArenaLogCount.Set(0)
	query = `SELECT COUNT(*) FROM log_arena_fights`
	err = e.connections.Characters.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	metrics.ArenaLogCount.Set(float64(count))

	// IP action logs (auth database)
	metrics.IPActionLogCount.Set(0)
	query = `SELECT COUNT(*) FROM logs_ip_actions`
	err = e.connections.Auth.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	metrics.IPActionLogCount.Set(float64(count))

	return nil
}

func (e *Exporter) collectInstanceMetrics() error {
	// Active instances
	metrics.ActiveInstanceCount.Set(0)
	query := `SELECT COUNT(*) FROM instance WHERE resettime > UNIX_TIMESTAMP()`
	var count int
	err := e.connections.Characters.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	metrics.ActiveInstanceCount.Set(float64(count))

	// Instances by difficulty
	metrics.InstancesByDifficulty.Reset()
	query = `SELECT difficulty, COUNT(*) FROM instance GROUP BY difficulty`
	rows, err := e.connections.Characters.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)
	for rows.Next() {
		var difficulty, count int
		if err := rows.Scan(&difficulty, &count); err != nil {
			return err
		}
		difficultyName := constants.GetDifficultyName(difficulty)
		metrics.InstancesByDifficulty.WithLabelValues(difficultyName).Set(float64(count))
	}

	// Completed encounters
	metrics.CompletedEncounters.Reset()
	query = `SELECT id, completedEncounters FROM instance WHERE completedEncounters > 0`
	rows, err = e.connections.Characters.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)
	for rows.Next() {
		var instanceID, encounters int
		if err := rows.Scan(&instanceID, &encounters); err != nil {
			return err
		}
		metrics.CompletedEncounters.WithLabelValues(fmt.Sprintf("%d", instanceID)).Set(float64(encounters))
	}

	// Instance resets
	metrics.InstanceResets.Reset()
	query = `SELECT mapid, difficulty, resettime FROM instance_reset`
	rows, err = e.connections.Characters.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)
	for rows.Next() {
		var mapID, difficulty, resetTime int
		if err := rows.Scan(&mapID, &difficulty, &resetTime); err != nil {
			return err
		}
		difficultyName := constants.GetDifficultyName(difficulty)
		metrics.InstanceResets.WithLabelValues(fmt.Sprintf("%d", mapID), difficultyName).Set(float64(resetTime))
	}

	// Characters in instances
	metrics.CharactersInInstances.Set(0)
	query = `SELECT COUNT(DISTINCT guid) FROM character_instance`
	err = e.connections.Characters.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	metrics.CharactersInInstances.Set(float64(count))

	// LFG data
	metrics.LFGDataCount.Reset()
	query = `SELECT state, COUNT(*) FROM lfg_data GROUP BY state`
	rows, err = e.connections.Characters.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)
	for rows.Next() {
		var state, count int
		if err := rows.Scan(&state, &count); err != nil {
			return err
		}
		stateName := constants.GetLFGStateName(state)
		metrics.LFGDataCount.WithLabelValues(stateName).Set(float64(count))
	}

	// Lag reports
	metrics.LagReportsCount.Set(0)
	query = `SELECT COUNT(*) FROM lag_reports`
	err = e.connections.Characters.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	metrics.LagReportsCount.Set(float64(count))

	// Instance saves
	metrics.InstanceSavesCount.Set(0)
	query = `SELECT COUNT(*) FROM instance_saved_go_state_data`
	err = e.connections.Characters.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}
	metrics.InstanceSavesCount.Set(float64(count))

	return nil
}

func (e *Exporter) collectNetworkMetrics() error {
	// Player latency statistics
	metrics.PlayerLatencyStats.Reset()

	// Average latency
	var avgLatency sql.NullFloat64
	query := `
		SELECT AVG(latency) 
		FROM characters 
		WHERE online = 1 
		AND latency > 0 
		AND (deleteDate IS NULL OR deleteDate = 0)
		AND (creation_date IS NULL OR creation_date < DATE_SUB(NOW(), INTERVAL 24 HOUR)) -- Exclude characters created in last 24 hours
		AND name NOT LIKE '%test%'
		AND name NOT LIKE '%admin%'
		AND name NOT LIKE '%gm%'
		AND name NOT LIKE '%dev%'
	`
	err := e.connections.Characters.QueryRow(query).Scan(&avgLatency)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err != sql.ErrNoRows && avgLatency.Valid {
		metrics.PlayerLatencyStats.WithLabelValues("average").Set(avgLatency.Float64)
		metrics.AverageLatency.Set(avgLatency.Float64)
	}

	// High latency players (>200ms)
	var highLatencyCount int
	query = `
		SELECT COUNT(*) 
		FROM characters 
		WHERE online = 1 
		AND latency > 200 
		AND (deleteDate IS NULL OR deleteDate = 0)
		AND (creation_date IS NULL OR creation_date < DATE_SUB(NOW(), INTERVAL 24 HOUR)) -- Exclude characters created in last 24 hours
		AND name NOT LIKE '%test%'
		AND name NOT LIKE '%admin%'
		AND name NOT LIKE '%gm%'
		AND name NOT LIKE '%dev%'
	`
	err = e.connections.Characters.QueryRow(query).Scan(&highLatencyCount)
	if err != nil {
		return err
	}
	metrics.HighLatencyPlayers.Set(float64(highLatencyCount))
	metrics.PlayerLatencyStats.WithLabelValues("high_latency").Set(float64(highLatencyCount))

	// Min/Max latency
	var minLatency, maxLatency sql.NullInt64
	query = `
		SELECT MIN(latency), MAX(latency) 
		FROM characters 
		WHERE online = 1 
		AND latency > 0 
		AND (deleteDate IS NULL OR deleteDate = 0)
		AND (creation_date IS NULL OR creation_date < DATE_SUB(NOW(), INTERVAL 24 HOUR)) -- Exclude characters created in last 24 hours
		AND name NOT LIKE '%test%'
		AND name NOT LIKE '%admin%'
		AND name NOT LIKE '%gm%'
		AND name NOT LIKE '%dev%'
	`
	err = e.connections.Characters.QueryRow(query).Scan(&minLatency, &maxLatency)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err != sql.ErrNoRows {
		if minLatency.Valid {
			metrics.PlayerLatencyStats.WithLabelValues("min").Set(float64(minLatency.Int64))
		}
		if maxLatency.Valid {
			metrics.PlayerLatencyStats.WithLabelValues("max").Set(float64(maxLatency.Int64))
		}
	}

	// IP bans
	metrics.IPBannedCount.Set(0)
	var bannedCount int
	query = `SELECT COUNT(*) FROM ip_banned`
	err = e.connections.Auth.QueryRow(query).Scan(&bannedCount)
	if err != nil {
		return err
	}
	metrics.IPBannedCount.Set(float64(bannedCount))

	// IP action logs by type
	metrics.IPActionLogsByType.Reset()
	query = `SELECT type, COUNT(*) FROM logs_ip_actions GROUP BY type`
	rows, err := e.connections.Auth.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)
	for rows.Next() {
		var actionType, count int
		if err := rows.Scan(&actionType, &count); err != nil {
			return err
		}
		typeName := constants.GetIPActionTypeName(actionType)
		metrics.IPActionLogsByType.WithLabelValues(typeName).Set(float64(count))
	}

	// Lag reports by type
	metrics.LagReportsByType.Reset()
	query = `SELECT lagType, COUNT(*) FROM lag_reports GROUP BY lagType`
	rows, err = e.connections.Characters.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)
	for rows.Next() {
		var lagType, count int
		if err := rows.Scan(&lagType, &count); err != nil {
			return err
		}
		lagTypeName := constants.GetLagTypeName(lagType)
		metrics.LagReportsByType.WithLabelValues(lagTypeName).Set(float64(count))
	}

	// Network activity by IP (top 10 most active IPs)
	metrics.NetworkActivityByIP.Reset()
	query = `SELECT ip, COUNT(*) as activity FROM logs_ip_actions GROUP BY ip ORDER BY activity DESC LIMIT 10`
	rows, err = e.connections.Auth.Query(query)
	if err != nil {
		return err
	}
	defer database.CloseRowsWithLog(rows)
	for rows.Next() {
		var ip string
		var activity int
		if err := rows.Scan(&ip, &activity); err != nil {
			return err
		}
		metrics.NetworkActivityByIP.WithLabelValues(ip).Set(float64(activity))
	}

	return nil
}

func (e *Exporter) collectBattlegroundMetrics() error {
	// Reset vector metrics
	metrics.BattlegroundDesertersByType.Reset()
	metrics.BattlegroundStats.Reset()
	metrics.BattlegroundsByType.Reset()
	metrics.BattlegroundsByBracket.Reset()
	metrics.BattlegroundWinsByFaction.Reset()
	metrics.BattlegroundPlayerStats.Reset()
	metrics.BattlegroundTemplates.Reset()
	metrics.BattlegroundTemplateDetails.Reset()
	metrics.RecentBattlegrounds.Reset()

	// Battleground deserters
	var deserterCount int
	err := e.connections.Characters.QueryRow("SELECT COUNT(*) FROM battleground_deserters").Scan(&deserterCount)
	if err != nil {
		return fmt.Errorf("error querying battleground deserters: %v", err)
	}
	metrics.BattlegroundDeserters.Set(float64(deserterCount))

	// Battleground deserters by type
	rows, err := e.connections.Characters.Query(`
		SELECT type, COUNT(*) as count 
		FROM battleground_deserters 
		GROUP BY type
	`)
	if err != nil {
		return fmt.Errorf("error querying battleground deserters by type: %v", err)
	}
	defer database.CloseRowsWithLog(rows)

	for rows.Next() {
		var desertionType int
		var count int
		if err := rows.Scan(&desertionType, &count); err != nil {
			continue
		}
		metrics.BattlegroundDesertersByType.WithLabelValues(constants.GetDesertionTypeName(desertionType)).Set(float64(count))
	}

	// Random battleground queue
	var queueCount int
	err = e.connections.Characters.QueryRow("SELECT COUNT(*) FROM character_battleground_random").Scan(&queueCount)
	if err != nil {
		return fmt.Errorf("error querying random battleground queue: %v", err)
	}
	metrics.RandomBattlegroundQueue.Set(float64(queueCount))

	// Battleground statistics
	var totalBattlegrounds, totalPlayers int
	err = e.connections.Characters.QueryRow(`
		SELECT 
			COUNT(DISTINCT id) as total_battlegrounds,
			COUNT(DISTINCT character_guid) as total_players
		FROM pvpstats_battlegrounds bg
		LEFT JOIN pvpstats_players bp ON bg.id = bp.battleground_id
	`).Scan(&totalBattlegrounds, &totalPlayers)
	if err != nil {
		return fmt.Errorf("error querying battleground stats: %v", err)
	}
	metrics.BattlegroundStats.WithLabelValues("total_battlegrounds").Set(float64(totalBattlegrounds))
	metrics.BattlegroundStats.WithLabelValues("total_players").Set(float64(totalPlayers))

	// Battlegrounds by type
	rows, err = e.connections.Characters.Query(`
		SELECT type, COUNT(*) as count 
		FROM pvpstats_battlegrounds 
		GROUP BY type
	`)
	if err != nil {
		return fmt.Errorf("error querying battlegrounds by type: %v", err)
	}
	defer database.CloseRowsWithLog(rows)

	for rows.Next() {
		var bgType int
		var count int
		if err := rows.Scan(&bgType, &count); err != nil {
			continue
		}
		metrics.BattlegroundsByType.WithLabelValues(constants.GetBattlegroundTypeName(bgType)).Set(float64(count))
	}

	// Battlegrounds by bracket
	rows, err = e.connections.Characters.Query(`
		SELECT bracket_id, COUNT(*) as count 
		FROM pvpstats_battlegrounds 
		GROUP BY bracket_id
	`)
	if err != nil {
		return fmt.Errorf("error querying battlegrounds by bracket: %v", err)
	}
	defer database.CloseRowsWithLog(rows)

	for rows.Next() {
		var bracket int
		var count int
		if err := rows.Scan(&bracket, &count); err != nil {
			continue
		}
		metrics.BattlegroundsByBracket.WithLabelValues(fmt.Sprintf("bracket_%d", bracket)).Set(float64(count))
	}

	// Battleground wins by faction
	rows, err = e.connections.Characters.Query(`
		SELECT winner_faction, COUNT(*) as count 
		FROM pvpstats_battlegrounds 
		WHERE winner_faction IN (0, 1)
		GROUP BY winner_faction
	`)
	if err != nil {
		return fmt.Errorf("error querying battleground wins by faction: %v", err)
	}
	defer database.CloseRowsWithLog(rows)

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
		metrics.BattlegroundWinsByFaction.WithLabelValues(factionName).Set(float64(count))
	}

	// Battleground player statistics
	rows, err = e.connections.Characters.Query(`
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
	defer database.CloseRowsWithLog(rows)

	if rows.Next() {
		var totalParticipants, totalWinners int
		var avgKillingBlows, avgDeaths, avgHonorableKills, avgBonusHonor, avgDamageDone, avgHealingDone sql.NullFloat64
		if err := rows.Scan(&totalParticipants, &totalWinners, &avgKillingBlows, &avgDeaths, &avgHonorableKills, &avgBonusHonor, &avgDamageDone, &avgHealingDone); err == nil {
			metrics.BattlegroundPlayerStats.WithLabelValues("total_participants").Set(float64(totalParticipants))
			metrics.BattlegroundPlayerStats.WithLabelValues("total_winners").Set(float64(totalWinners))
			if avgKillingBlows.Valid {
				metrics.BattlegroundPlayerStats.WithLabelValues("avg_killing_blows").Set(avgKillingBlows.Float64)
			}
			if avgDeaths.Valid {
				metrics.BattlegroundPlayerStats.WithLabelValues("avg_deaths").Set(avgDeaths.Float64)
			}
			if avgHonorableKills.Valid {
				metrics.BattlegroundPlayerStats.WithLabelValues("avg_honorable_kills").Set(avgHonorableKills.Float64)
			}
			if avgBonusHonor.Valid {
				metrics.BattlegroundPlayerStats.WithLabelValues("avg_bonus_honor").Set(avgBonusHonor.Float64)
			}
			if avgDamageDone.Valid {
				metrics.BattlegroundPlayerStats.WithLabelValues("avg_damage_done").Set(avgDamageDone.Float64)
			}
			if avgHealingDone.Valid {
				metrics.BattlegroundPlayerStats.WithLabelValues("avg_healing_done").Set(avgHealingDone.Float64)
			}
		}
	}

	// Battleground templates
	rows, err = e.connections.World.Query(`
		SELECT ID, ScriptName, Comment, MinPlayersPerTeam, MaxPlayersPerTeam, MinLvl, MaxLvl, Weight
		FROM battleground_template
	`)
	if err != nil {
		return fmt.Errorf("error querying battleground templates: %v", err)
	}
	defer database.CloseRowsWithLog(rows)

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

		metrics.BattlegroundTemplates.WithLabelValues(fmt.Sprintf("%d", id), label).Set(float64(weight))

		// Add detailed template information
		metrics.BattlegroundTemplateDetails.WithLabelValues(
			fmt.Sprintf("%d", id),
			label,
			fmt.Sprintf("%d", minLvl),
			fmt.Sprintf("%d", maxLvl),
			fmt.Sprintf("%d", minPlayers),
			fmt.Sprintf("%d", maxPlayers),
		).Set(float64(weight))
	}

	// Recent battleground activity (last 24 hours, 7 days, 30 days)
	rows, err = e.connections.Characters.Query(`
		SELECT 
			SUM(CASE WHEN date >= DATE_SUB(NOW(), INTERVAL 24 HOUR) THEN 1 ELSE 0 END) as last_24h,
			SUM(CASE WHEN date >= DATE_SUB(NOW(), INTERVAL 7 DAY) THEN 1 ELSE 0 END) as last_7d,
			SUM(CASE WHEN date >= DATE_SUB(NOW(), INTERVAL 30 DAY) THEN 1 ELSE 0 END) as last_30d
		FROM pvpstats_battlegrounds
	`)
	if err != nil {
		return fmt.Errorf("error querying recent battleground activity: %v", err)
	}
	defer database.CloseRowsWithLog(rows)

	if rows.Next() {
		var last24h, last7d, last30d sql.NullInt64
		if err := rows.Scan(&last24h, &last7d, &last30d); err == nil {
			if last24h.Valid {
				metrics.RecentBattlegrounds.WithLabelValues("last_24h").Set(float64(last24h.Int64))
			}
			if last7d.Valid {
				metrics.RecentBattlegrounds.WithLabelValues("last_7d").Set(float64(last7d.Int64))
			}
			if last30d.Valid {
				metrics.RecentBattlegrounds.WithLabelValues("last_30d").Set(float64(last30d.Int64))
			}
		}
	}

	return nil
}
