# WoW Private Server Prometheus Exporter

A comprehensive Prometheus exporter for World of Warcraft private servers running on AzerothCore. This exporter collects metrics from your WoW server databases and exposes them for monitoring with Prometheus and Grafana.

## Features

### 🎮 Player Metrics
- Players online by faction (Alliance/Horde)
- Total players by faction
- Players by level and class
- Max-level characters by faction

### 📧 Mail System
- Total mail messages
- Mail by faction
- Mail with items
- Unread mail count

### 👥 Account Management
- Total accounts
- Online accounts
- Banned accounts
- GM accounts

### ⚔️ Battleground & PvP
- Battleground templates and configuration
- Random battleground queue
- Battleground statistics
- Battleground deserters
- Player performance in battlegrounds
- Recent battleground activity

### 🏰 Guilds & Social
- Total guild count
- Guild events
- Chat channels and bans

### 🏪 Economy
- Auction house activity by faction
- Money transaction logs

### 🏛️ Instances & Raids
- Active instances
- Instances by difficulty
- Completed encounters
- Instance resets
- Characters in instances
- LFG data

### 🌐 Network & Performance
- Player latency statistics
- Average latency
- High latency players
- IP bans and actions
- Lag reports
- Network activity by IP

## Quick Start

### Option 1: Nix/NixOS (Recommended)

#### Development Environment
```bash
# Enter development shell
nix develop

# Or with legacy nix-shell
nix-shell
```

#### NixOS Service
Add to your `configuration.nix`:
```nix
{ config, lib, pkgs, ... }:

let
  wowExporterModule = (import ./flake.nix).nixosModules.default;
in
{
  imports = [ wowExporterModule ];

  services.wow-private-server-exporter = {
    enable = true;
    port = 7000;
    database = {
      user = "wow_user";
      password = "your_password";
      host = "your-database-host";
      port = 3306;
    };
  };
}
```

### Option 2: Manual Setup

#### 1. Environment Setup

Set your database connection details:

```bash
export WOW_DB_USER=root
export WOW_DB_PASS=password
export WOW_DB_HOST=your-server-ip
export WOW_DB_PORT=3306  # Optional, defaults to 3306
export PORT=7000         # Optional, defaults to 7000
```

#### 2. Run the Exporter

```bash
go run main.go
```

Or build and run:

```bash
go build -o prom-azerothcore-exporter main.go
./prom-azerothcore-exporter
```

#### 3. Verify Metrics

Visit `http://localhost:7000/metrics` to see the raw Prometheus metrics.

## Grafana Dashboard Setup

### 1. Import the Dashboard

1. Open Grafana
2. Go to **Dashboards** → **Import**
3. Upload the `grafana-dashboard.json` file
4. Select your Prometheus data source
5. Click **Import**

### 2. Configure Prometheus Data Source

Make sure your Prometheus is configured to scrape the exporter:

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'wow-private-server'
    static_configs:
      - targets: ['localhost:7000']
    scrape_interval: 15s
```

## Dashboard Panels

### 📊 Overview Panels
- **Players Online by Faction** - Real-time player count by Alliance/Horde
- **Players by Class** - Class distribution over time
- **Key Statistics** - Total accounts, online accounts, guilds, mail

### ⚔️ Battleground Monitoring
- **Battleground Templates** - Available battleground types
- **Random Battleground Queue** - Players waiting for random BGs
- **Battleground Statistics** - Overall BG activity
- **Battleground Deserters** - Players leaving BGs

### 🏛️ Instance & Raid Tracking
- **Active Instances** - Currently running instances
- **Auction House Activity** - Economic activity by faction

### 🌐 Network Performance
- **Average Player Latency** - Server performance
- **High Latency Players** - Players experiencing issues

## Example Prometheus Queries

### Player Activity
```promql
# Players online by faction
wow_players_online

# Most popular class
topk(5, sum by (class) (wow_players_by_class))

# Players by level range
sum by (level) (wow_players_by_level{level=~"8[0-9]"})
```

### Battleground Activity
```promql
# Random BG queue
wow_random_battleground_queue

# Battleground templates
wow_battleground_templates

# Recent BG activity
wow_recent_battlegrounds{time_period="last_24h"}
```

### Server Health
```promql
# Average latency
wow_average_latency_ms

# High latency players
wow_high_latency_players

# Server uptime
wow_server_uptime_seconds
```

## Database Requirements

The exporter connects to three AzerothCore databases:

- **acore_characters** - Player data, mail, guilds, instances
- **acore_auth** - Account data, bans, IP actions
- **acore_world** - Battleground templates, game data

## Configuration

### NixOS Module Configuration

The NixOS module provides a declarative way to configure the exporter:

```nix
services.wow-private-server-exporter = {
  enable = true;
  port = 7000;
  
  # Database configuration
  database = {
    user = "wow_user";
    password = "your_password";
    host = "your-database-host";
    port = 3306;
  };
  
  # Optional: Use environment file for sensitive data
  environmentFile = "/etc/wow-exporter/env";
  
  # Optional: Custom user/group
  user = "wow-exporter";
  group = "wow-exporter";
};
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `WOW_DB_USER` | - | Database username |
| `WOW_DB_PASS` | - | Database password |
| `WOW_DB_HOST` | - | Database host |
| `WOW_DB_PORT` | 3306 | Database port |
| `WOW_DB_DSN` | - | Full DSN (overrides individual vars) |
| `PORT` | 7000 | Exporter port |

### Full DSN Example
```bash
export WOW_DB_DSN="user:pass@tcp(host:3306)/acore_characters?parseTime=true"
```

### Security with Environment Files

For production deployments, use environment files to keep sensitive data out of the Nix store:

```nix
# Create secure environment file
environment.etc."wow-exporter/env" = {
  text = ''
    WOW_DB_USER=wow_user
    WOW_DB_PASS=your_secure_password
    WOW_DB_HOST=your-database-host
    WOW_DB_PORT=3306
  '';
  mode = "0600";
  user = "wow-exporter";
  group = "wow-exporter";
};

# Use environment file in service
services.wow-private-server-exporter = {
  enable = true;
  environmentFile = "/etc/wow-exporter/env";
  # ... other config
};
```

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Check database host and port
   - Verify database credentials
   - Ensure database is running

2. **Table Not Found**
   - Verify database names (acore_characters, acore_auth, acore_world)
   - Check if tables exist in the correct databases

3. **Empty Metrics**
   - Some metrics may be empty if there's no activity
   - Battleground metrics require PvP activity
   - Instance metrics require dungeon/raid activity

### Debug Mode

Run with verbose logging:
```bash
go run main.go 2>&1 | tee exporter.log
```

## Metrics Reference

### Player Metrics
- `wow_players_online{faction}` - Players currently online
- `wow_players_total{faction}` - Total players by faction
- `wow_players_by_level{level,faction}` - Players by level
- `wow_players_by_class{class,faction}` - Players by class

### Battleground Metrics
- `wow_battleground_templates{template_id,script_name}` - BG templates
- `wow_random_battleground_queue` - Players in random BG queue
- `wow_battleground_stats{stat}` - BG statistics
- `wow_battleground_deserters` - BG deserters

### Server Metrics
- `wow_server_uptime_seconds` - Server uptime
- `wow_average_latency_ms` - Average player latency
- `wow_high_latency_players` - Players with high latency

## Contributing

Feel free to submit issues and enhancement requests!

## License

This project is open source and available under the MIT License. 