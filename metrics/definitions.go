package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Player metrics
var (
	PlayersOnline = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_players_online",
			Help: "Number of players currently online",
		},
		[]string{"faction"},
	)

	PlayersTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_players_total",
			Help: "Total number of players",
		},
		[]string{"faction"},
	)

	PlayersByLevel = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_players_by_level",
			Help: "Number of players by level",
		},
		[]string{"level", "faction"},
	)

	PlayersByClass = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_players_by_class",
			Help: "Number of players by class",
		},
		[]string{"class", "faction"},
	)
)

// Mail metrics
var (
	MailTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_mail_total",
			Help: "Total number of mail messages",
		},
	)

	MailByFaction = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_mail_by_faction",
			Help: "Number of mail messages by faction",
		},
		[]string{"faction"},
	)

	MailWithItems = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_mail_with_items",
			Help: "Number of mail messages with items",
		},
	)
)

// Account metrics
var (
	AccountsTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_accounts_total",
			Help: "Total number of accounts",
		},
	)

	AccountsOnline = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_accounts_online",
			Help: "Number of accounts currently online",
		},
	)

	AccountsBanned = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_accounts_banned",
			Help: "Number of banned accounts",
		},
	)
)

// Server metrics
var (
	ServerUptime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_server_uptime_seconds",
			Help: "Server uptime in seconds",
		},
	)

	ServerMaxPlayers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_server_max_players",
			Help: "Maximum number of players recorded",
		},
	)
)

// Auction metrics
var (
	AuctionCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_auction_count",
			Help: "Number of active auctions by house (faction)",
		},
		[]string{"house"},
	)
)

// Guild metrics
var (
	GuildCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_guild_count",
			Help: "Number of guilds",
		},
	)
)

// Character metrics
var (
	MaxLevelCharCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_max_level_characters",
			Help: "Number of max-level characters by faction",
		},
		[]string{"faction"},
	)

	UnreadMailCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_unread_mail_count",
			Help: "Number of unread mail messages",
		},
	)

	GMAccountCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_gm_account_count",
			Help: "Number of accounts with GM level",
		},
	)

	LastServerRestart = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_server_last_restart_timestamp",
			Help: "Timestamp of the last server restart (unix time)",
		},
	)

	BannedCharCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_banned_characters",
			Help: "Number of banned characters",
		},
	)
)

// Chat and activity metrics
var (
	ChannelCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_channel_count",
			Help: "Number of chat channels",
		},
	)

	ChannelBans = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_channel_bans",
			Help: "Number of channel bans",
		},
	)

	LogCountByType = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_log_count",
			Help: "Number of log entries by type",
		},
		[]string{"type"},
	)

	GuildEventCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_guild_events",
			Help: "Number of guild events",
		},
	)

	MoneyLogCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_money_logs",
			Help: "Number of money transaction logs",
		},
	)

	EncounterLogCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_encounter_logs",
			Help: "Number of encounter logs",
		},
	)

	ArenaLogCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_arena_logs",
			Help: "Number of arena fight logs",
		},
	)

	IPActionLogCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_ip_action_logs",
			Help: "Number of IP action logs",
		},
	)
)

// Instance and raid metrics
var (
	ActiveInstanceCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_active_instances",
			Help: "Number of active instances",
		},
	)

	InstancesByDifficulty = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_instances_by_difficulty",
			Help: "Number of instances by difficulty",
		},
		[]string{"difficulty"},
	)

	CompletedEncounters = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_completed_encounters",
			Help: "Number of completed encounters by instance",
		},
		[]string{"instance_id"},
	)

	InstanceResets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_instance_resets",
			Help: "Instance reset times by map and difficulty",
		},
		[]string{"map_id", "difficulty"},
	)

	CharactersInInstances = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_characters_in_instances",
			Help: "Number of characters currently in instances",
		},
	)

	LFGDataCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_lfg_data",
			Help: "Number of LFG entries by state",
		},
		[]string{"state"},
	)

	LagReportsCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_lag_reports",
			Help: "Number of lag reports",
		},
	)

	InstanceSavesCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_instance_saves",
			Help: "Number of saved instance states",
		},
	)
)

// Networking metrics
var (
	PlayerLatencyStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_player_latency",
			Help: "Player latency statistics",
		},
		[]string{"stat"},
	)

	IPBannedCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_ip_banned_count",
			Help: "Number of banned IP addresses",
		},
	)

	IPActionLogsByType = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_ip_action_logs_by_type",
			Help: "Number of IP action logs by type",
		},
		[]string{"type"},
	)

	LagReportsByType = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_lag_reports_by_type",
			Help: "Number of lag reports by type",
		},
		[]string{"lag_type"},
	)

	AverageLatency = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_average_latency_ms",
			Help: "Average player latency in milliseconds",
		},
	)

	HighLatencyPlayers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_high_latency_players",
			Help: "Number of players with high latency (>200ms)",
		},
	)

	NetworkActivityByIP = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_network_activity_by_ip",
			Help: "Network activity by IP address (top 10)",
		},
		[]string{"ip"},
	)
)

// Battleground metrics
var (
	BattlegroundDeserters = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_battleground_deserters",
			Help: "Number of battleground deserters",
		},
	)

	BattlegroundDesertersByType = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battleground_deserters_by_type",
			Help: "Number of battleground deserters by type",
		},
		[]string{"desertion_type"},
	)

	RandomBattlegroundQueue = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wow_random_battleground_queue",
			Help: "Number of players in random battleground queue",
		},
	)

	BattlegroundStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battleground_stats",
			Help: "Battleground statistics",
		},
		[]string{"stat"},
	)

	BattlegroundsByType = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battlegrounds_by_type",
			Help: "Number of battlegrounds by type",
		},
		[]string{"battleground_type"},
	)

	BattlegroundsByBracket = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battlegrounds_by_bracket",
			Help: "Number of battlegrounds by bracket",
		},
		[]string{"bracket"},
	)

	BattlegroundWinsByFaction = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battleground_wins_by_faction",
			Help: "Number of battleground wins by faction",
		},
		[]string{"faction"},
	)

	BattlegroundPlayerStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battleground_player_stats",
			Help: "Battleground player statistics",
		},
		[]string{"stat"},
	)

	BattlegroundTemplates = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battleground_templates",
			Help: "Battleground template information",
		},
		[]string{"template_id", "script_name"},
	)

	BattlegroundTemplateDetails = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_battleground_template_details",
			Help: "Detailed battleground template information",
		},
		[]string{"template_id", "name", "min_level", "max_level", "min_players", "max_players"},
	)

	RecentBattlegrounds = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wow_recent_battlegrounds",
			Help: "Recent battleground activity",
		},
		[]string{"time_period"},
	)
)
