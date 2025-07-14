package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/scottjab/prom-azerothcore-exporter/config"
	"github.com/scottjab/prom-azerothcore-exporter/internal/exporter"
	"github.com/scottjab/prom-azerothcore-exporter/pkg/database"
)

func main() {
	cfg := config.Load()

	connections, err := database.NewConnections(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to databases: %v", err)
	}
	defer connections.Close()

	exp := exporter.NewExporter(connections)
	defer exp.Close()

	prometheus.MustRegister(exp)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`
			<html>
			<head><title>WoW Private Server Exporter</title></head>
			<body>
				<h1>WoW Private Server Exporter</h1>
				<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>
		`))
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})

	log.Printf("Starting WoW Private Server Exporter on port %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, nil))
}
