package exporter

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/scottjab/prom-azerothcore-exporter/metrics"
	"github.com/scottjab/prom-azerothcore-exporter/pkg/database"
)

// Exporter implements the Prometheus Collector interface
type Exporter struct {
	connections *database.Connections
}

// NewExporter creates a new exporter instance
func NewExporter(connections *database.Connections) *Exporter {
	return &Exporter{
		connections: connections,
	}
}

// Close closes the exporter and its database connections
func (e *Exporter) Close() {
	if e.connections != nil {
		e.connections.Close()
	}
}

// Collect implements prometheus.Collector
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
	metrics.PlayersOnline.Collect(ch)
	metrics.PlayersTotal.Collect(ch)
	metrics.PlayersByLevel.Collect(ch)
	metrics.PlayersByClass.Collect(ch)
	metrics.OnlinePlayersByLevel.Collect(ch)
	metrics.MailTotal.Collect(ch)
	metrics.MailByFaction.Collect(ch)
	metrics.MailWithItems.Collect(ch)
	metrics.AccountsTotal.Collect(ch)
	metrics.AccountsOnline.Collect(ch)
	metrics.AccountsBanned.Collect(ch)
	metrics.ServerUptime.Collect(ch)
	metrics.ServerMaxPlayers.Collect(ch)
	metrics.AuctionCount.Collect(ch)
	metrics.GuildCount.Collect(ch)
	metrics.MaxLevelCharCount.Collect(ch)
	metrics.UnreadMailCount.Collect(ch)
	metrics.GMAccountCount.Collect(ch)
	metrics.LastServerRestart.Collect(ch)
	metrics.BannedCharCount.Collect(ch)
	metrics.ChannelCount.Collect(ch)
	metrics.ChannelBans.Collect(ch)
	metrics.LogCountByType.Collect(ch)
	metrics.GuildEventCount.Collect(ch)
	metrics.MoneyLogCount.Collect(ch)
	metrics.EncounterLogCount.Collect(ch)
	metrics.ArenaLogCount.Collect(ch)
	metrics.IPActionLogCount.Collect(ch)
	metrics.ActiveInstanceCount.Collect(ch)
	metrics.InstancesByDifficulty.Collect(ch)
	metrics.CompletedEncounters.Collect(ch)
	metrics.InstanceResets.Collect(ch)
	metrics.CharactersInInstances.Collect(ch)
	metrics.LFGDataCount.Collect(ch)
	metrics.LagReportsCount.Collect(ch)
	metrics.InstanceSavesCount.Collect(ch)
	metrics.PlayerLatencyStats.Collect(ch)
	metrics.IPBannedCount.Collect(ch)
	metrics.IPActionLogsByType.Collect(ch)
	metrics.LagReportsByType.Collect(ch)
	metrics.AverageLatency.Collect(ch)
	metrics.HighLatencyPlayers.Collect(ch)
	metrics.NetworkActivityByIP.Collect(ch)
	metrics.BattlegroundDeserters.Collect(ch)
	metrics.BattlegroundDesertersByType.Collect(ch)
	metrics.RandomBattlegroundQueue.Collect(ch)
	metrics.BattlegroundStats.Collect(ch)
	metrics.BattlegroundsByType.Collect(ch)
	metrics.BattlegroundsByBracket.Collect(ch)
	metrics.BattlegroundWinsByFaction.Collect(ch)
	metrics.BattlegroundPlayerStats.Collect(ch)
	metrics.BattlegroundTemplates.Collect(ch)
	metrics.BattlegroundTemplateDetails.Collect(ch)
	metrics.RecentBattlegrounds.Collect(ch)
	metrics.ActiveBattlegrounds.Collect(ch)
	metrics.ActiveBattlegroundPlayers.Collect(ch)
	metrics.ActiveBattlegroundTotal.Collect(ch)
}

// Describe implements prometheus.Collector
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	metrics.PlayersOnline.Describe(ch)
	metrics.PlayersTotal.Describe(ch)
	metrics.PlayersByLevel.Describe(ch)
	metrics.PlayersByClass.Describe(ch)
	metrics.OnlinePlayersByLevel.Describe(ch)
	metrics.MailTotal.Describe(ch)
	metrics.MailByFaction.Describe(ch)
	metrics.MailWithItems.Describe(ch)
	metrics.AccountsTotal.Describe(ch)
	metrics.AccountsOnline.Describe(ch)
	metrics.AccountsBanned.Describe(ch)
	metrics.ServerUptime.Describe(ch)
	metrics.ServerMaxPlayers.Describe(ch)
	metrics.AuctionCount.Describe(ch)
	metrics.GuildCount.Describe(ch)
	metrics.MaxLevelCharCount.Describe(ch)
	metrics.UnreadMailCount.Describe(ch)
	metrics.GMAccountCount.Describe(ch)
	metrics.LastServerRestart.Describe(ch)
	metrics.BannedCharCount.Describe(ch)
	metrics.ChannelCount.Describe(ch)
	metrics.ChannelBans.Describe(ch)
	metrics.LogCountByType.Describe(ch)
	metrics.GuildEventCount.Describe(ch)
	metrics.MoneyLogCount.Describe(ch)
	metrics.EncounterLogCount.Describe(ch)
	metrics.ArenaLogCount.Describe(ch)
	metrics.IPActionLogCount.Describe(ch)
	metrics.ActiveInstanceCount.Describe(ch)
	metrics.InstancesByDifficulty.Describe(ch)
	metrics.CompletedEncounters.Describe(ch)
	metrics.InstanceResets.Describe(ch)
	metrics.CharactersInInstances.Describe(ch)
	metrics.LFGDataCount.Describe(ch)
	metrics.LagReportsCount.Describe(ch)
	metrics.InstanceSavesCount.Describe(ch)
	metrics.PlayerLatencyStats.Describe(ch)
	metrics.IPBannedCount.Describe(ch)
	metrics.IPActionLogsByType.Describe(ch)
	metrics.LagReportsByType.Describe(ch)
	metrics.AverageLatency.Describe(ch)
	metrics.HighLatencyPlayers.Describe(ch)
	metrics.NetworkActivityByIP.Describe(ch)
	metrics.BattlegroundDeserters.Describe(ch)
	metrics.BattlegroundDesertersByType.Describe(ch)
	metrics.RandomBattlegroundQueue.Describe(ch)
	metrics.BattlegroundStats.Describe(ch)
	metrics.BattlegroundsByType.Describe(ch)
	metrics.BattlegroundsByBracket.Describe(ch)
	metrics.BattlegroundWinsByFaction.Describe(ch)
	metrics.BattlegroundPlayerStats.Describe(ch)
	metrics.BattlegroundTemplates.Describe(ch)
	metrics.BattlegroundTemplateDetails.Describe(ch)
	metrics.RecentBattlegrounds.Describe(ch)
	metrics.ActiveBattlegrounds.Describe(ch)
	metrics.ActiveBattlegroundPlayers.Describe(ch)
	metrics.ActiveBattlegroundTotal.Describe(ch)
}

// Helper function for writing HTTP responses
func WriteWithLog(w http.ResponseWriter, data []byte) {
	if _, err := w.Write(data); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
