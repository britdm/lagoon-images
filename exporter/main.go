package main

import (
	"log"
	"net/http"
	"strings"

	ps "github.com/mitchellh/go-ps"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getCurrentProccessCount() float64 {
	processList, err := ps.Processes()
	if err != nil {
		log.Println("Failed, to list current processes.")
	}

	count := -1

	for pid := range processList {
		var process ps.Process
		process = processList[pid]
		if strings.Contains(process.Executable(), "php-fpm") {
			count++
		}
	}
	return float64(count)
}

var (
	processCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "phpfpm_processes_count",
			Help: "Number of phpfpm process running.",
		})
)

func processCounter() {
	go func() {
		for {
			processCount.Set(getCurrentProccessCount())
		}
	}()
}

func main() {
	http.Handle("/pod-metrics", promhttp.Handler())

	processCounter()

	http.ListenAndServe(":8082", nil)
}
