{
  description = "WoW Private Server Prometheus Exporter";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "prom-azerothcore-exporter";
          version = "1.0.0";
          src = ./.;

          vendorHash = null;

          meta = with pkgs.lib; {
            description = "Prometheus exporter for WoW private servers running AzerothCore";
            homepage = "https://github.com/scottjab/prom-azerothcore-exporter";
            license = licenses.mit;
            maintainers = [ ];
            platforms = platforms.unix;
          };
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/prom-azerothcore-exporter";
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gcc
            mysql80
          ];

          shellHook = ''
            echo "WoW Private Server Exporter Development Environment"
            echo "Available commands:"
            echo "  go run .                    - Run the exporter"
            echo "  go build -o exporter .      - Build the exporter"
            echo "  go test ./...               - Run tests"
            echo ""
            echo "Environment variables:"
            echo "  WOW_DB_USER                 - Database username"
            echo "  WOW_DB_PASS                 - Database password"
            echo "  WOW_DB_HOST                 - Database host"
            echo "  WOW_DB_PORT                 - Database port (default: 3306)"
            echo "  PORT                        - Exporter port (default: 7000)"
          '';
        };
      }
    ) // {
      nixosModules.default = { config, lib, pkgs, ... }:
        with lib;
        let
          cfg = config.services.wow-private-server-exporter;
        in
        {
          options.services.wow-private-server-exporter = {
            enable = mkEnableOption "WoW Private Server Prometheus Exporter";

            package = mkOption {
              type = types.package;
              default = self.packages.${pkgs.system}.default;
              description = "The WoW Private Server Exporter package to use.";
            };

            port = mkOption {
              type = types.port;
              default = 7000;
              description = "Port to run the exporter on.";
            };

            database = {
              user = mkOption {
                type = types.str;
                description = "Database username.";
              };

              password = mkOption {
                type = types.str;
                description = "Database password.";
              };

              host = mkOption {
                type = types.str;
                description = "Database host.";
              };

              port = mkOption {
                type = types.port;
                default = 3306;
                description = "Database port.";
              };
            };

            environmentFile = mkOption {
              type = types.nullOr types.path;
              default = null;
              description = "Environment file to load database credentials from.";
            };

            user = mkOption {
              type = types.str;
              default = "wow-exporter";
              description = "User account under which the exporter runs.";
            };

            group = mkOption {
              type = types.str;
              default = "wow-exporter";
              description = "Group under which the exporter runs.";
            };
          };

          config = mkIf cfg.enable {
            users.users = mkIf (cfg.user == "wow-exporter") {
              wow-exporter = {
                isSystemUser = true;
                group = cfg.group;
                description = "WoW Private Server Exporter user";
                home = "/var/lib/wow-exporter";
                createHome = true;
              };
            };

            users.groups = mkIf (cfg.group == "wow-exporter") {
              wow-exporter = { };
            };

            systemd.services.wow-private-server-exporter = {
              description = "WoW Private Server Prometheus Exporter";
              wantedBy = [ "multi-user.target" ];
              after = [ "network.target" ];

              serviceConfig = {
                Type = "simple";
                User = cfg.user;
                Group = cfg.group;
                ExecStart = "${cfg.package}/bin/prom-azerothcore-exporter";
                Restart = "always";
                RestartSec = "10";
                ProtectSystem = "strict";
                ProtectHome = true;
                NoNewPrivileges = true;
                PrivateTmp = true;
                PrivateDevices = true;
                ProtectKernelTunables = true;
                ProtectKernelModules = true;
                ProtectControlGroups = true;
                RestrictRealtime = true;
                RestrictSUIDSGID = true;
                ReadWritePaths = [ "/var/lib/wow-exporter" ];
                Environment = [
                  "WOW_DB_USER=${cfg.database.user}"
                  "WOW_DB_PASS=${cfg.database.password}"
                  "WOW_DB_HOST=${cfg.database.host}"
                  "WOW_DB_PORT=${toString cfg.database.port}"
                  "PORT=${toString cfg.port}"
                ];
              } // (if cfg.environmentFile != null then {
                EnvironmentFile = [ cfg.environmentFile ];
              } else { });
            };

            # Prometheus configuration
            services.prometheus.scrapeConfigs = mkIf config.services.prometheus.enable [
              {
                job_name = "wow-private-server";
                static_configs = [{
                  targets = [ "localhost:${toString cfg.port}" ];
                }];
                scrape_interval = "15s";
                metrics_path = "/metrics";
              }
            ];
          };
        };
    };
} 