# Grafana Dashboards for AzerothCore Prometheus Exporter

This directory contains comprehensive Grafana dashboards for monitoring your AzerothCore server using the Prometheus exporter.

## Available Dashboards

### 1. `grafana-dashboard.json` - Basic Dashboard
A simplified dashboard with essential metrics for quick server monitoring.

### 2. `grafana-dashboard-comprehensive.json` - Comprehensive Dashboard
A full-featured dashboard with all available metrics organized into logical sections.

## Dashboard Features

### Server Status Overview
- **Online Players**: Real-time player count by faction
- **Server Uptime**: Current server uptime in seconds
- **Max Players Recorded**: Peak player count
- **Last Restart**: Timestamp of the last server restart

### Player Demographics
- **Players by Level**: Distribution of players across levels
- **Players by Class**: Class distribution with faction breakdown

### Account Management
- **Account Statistics**: Total, online, banned, and GM accounts
- **Character Statistics**: Max-level characters, banned characters, guild count

### Mail System
- **Mail Overview**: Total mail, mail with items, unread mail
- **Mail by Faction**: Mail distribution across factions

### Economy & Trading
- **Auction House**: Active auctions by house (faction)
- **Money Transactions**: Money transaction logs

### Network & Performance
- **Latency Overview**: Average latency and high latency players
- **Network Security**: Banned IPs and lag reports

### Instances & Raids
- **Instance Activity**: Active instances, players in instances, instance saves
- **Instances by Difficulty**: Instance distribution by difficulty
- **LFG Activity**: Looking for Group data by state

### Battlegrounds
- **BG Activity**: Random battleground queue and deserters
- **BG by Type**: Battleground distribution by type
- **BG Wins by Faction**: Battleground wins by faction

### System Logs & Events
- **Log Activity**: System logs by type
- **Event Counts**: Guild events, encounter logs, arena logs, IP action logs

### Chat & Communication
- **Chat Channels**: Total channels and channel bans

## Installation Instructions

### Prerequisites
1. Grafana instance running
2. Prometheus configured to scrape your AzerothCore exporter
3. AzerothCore Prometheus exporter running and collecting metrics

### Import Steps

1. **Open Grafana**
   - Navigate to your Grafana instance
   - Log in with appropriate credentials

2. **Import Dashboard**
   - Click the "+" icon in the sidebar
   - Select "Import"
   - Click "Upload JSON file" or paste the JSON content
   - Select the dashboard file you want to import

3. **Configure Data Source**
   - Select your Prometheus data source
   - Click "Import"

4. **Verify Metrics**
   - Check that all panels are displaying data
   - If panels show "No data", verify:
     - Prometheus is scraping your exporter
     - Exporter is running and collecting metrics
     - Metric names match those in your exporter

## Configuration

### Data Source
Ensure your Prometheus data source is configured with:
- **URL**: Your Prometheus server URL (e.g., `http://localhost:9090`)
- **Access**: Server (default) or Browser (if using Grafana Cloud)

### Variables
The comprehensive dashboard includes a `faction` variable that allows filtering by:
- Alliance
- Horde
- All factions

### Refresh Rate
Default refresh rate is set to 30 seconds. You can adjust this in the dashboard settings.

## Customization

### Adding New Panels
To add new panels for additional metrics:

1. **Edit Dashboard**
   - Click the dashboard title
   - Select "Settings" → "JSON Model"

2. **Add Panel**
   ```json
   {
     "id": 34,
     "title": "New Metric",
     "type": "stat",
     "targets": [
       {
         "expr": "your_metric_name",
         "legendFormat": "{{label}}"
       }
     ],
     "gridPos": {"h": 8, "w": 12, "x": 0, "y": 90}
   }
   ```

3. **Save Changes**
   - Click "Save changes" to update the dashboard

### Modifying Existing Panels
- Click on any panel title
- Select "Edit" to modify queries, visualization, or layout
- Use the query editor to modify PromQL expressions

## Troubleshooting

### Common Issues

1. **No Data in Panels**
   - Check Prometheus targets page
   - Verify exporter is running: `curl http://exporter:port/metrics`
   - Check metric names match exactly

2. **Incorrect Data**
   - Verify PromQL queries are correct
   - Check for metric label mismatches
   - Ensure time range is appropriate

3. **Performance Issues**
   - Reduce refresh rate for large datasets
   - Use time range selectors to limit data
   - Consider using recording rules for complex queries

### Useful PromQL Queries

```promql
# Total online players
sum(wow_players_online)

# Players by faction
wow_players_online

# Server uptime in hours
wow_server_uptime_seconds / 3600

# Average latency
wow_average_latency_ms

# Active instances
wow_active_instances
```

## Support

For issues with:
- **Dashboard**: Check Grafana logs and documentation
- **Metrics**: Verify exporter configuration and database connectivity
- **Prometheus**: Check Prometheus configuration and targets

## Contributing

To improve these dashboards:
1. Export your modified dashboard as JSON
2. Update the appropriate dashboard file
3. Update this README with new features
4. Test with different server configurations

## License

These dashboards are provided as-is for use with the AzerothCore Prometheus exporter. 