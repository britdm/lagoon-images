package main

import (
	"log"
	"net/http"
	"strings"

	ps "github.com/mitchellh/go-ps"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type processManager struct {
}

func (pm processManager) getCurrentProccessCount() float64 {
	log.Println("called process")
	processList, err := ps.Processes()
	if err != nil {
		log.Println("Failed to list current processes.")
	}

	count := -1

	phpName := "php-fpm"
	for pid := range processList {
		process := processList[pid]
		// log.Printf("process: %v", process.Executable())
		if strings.Contains(process.Executable(), phpName) {
			count++
			log.Println(count)
		}
	}
	return float64(count)
}

// Descriptors used by the ClusterManagerCollector below.
var (
	numProcesses = prometheus.NewDesc(
		"phpfpm_processes_count",
		"Number of phpfpm process running.",
		[]string{}, nil,
	)
)

// Describe is implemented with DescribeByCollect. That's possible because the
// Collect method will always return the same two metrics with the same two
// descriptors.
func (ps processManager) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(ps, ch)
}

// Collect first triggers the ReallyExpensiveAssessmentOfTheSystemState. Then it
// creates constant metrics for each host on the fly based on the returned data.
//
// Note that Collect could be called concurrently, so we depend on
// ReallyExpensiveAssessmentOfTheSystemState to be concurrency-safe.
func (ps processManager) Collect(ch chan<- prometheus.Metric) {
	processCount := ps.getCurrentProccessCount()
	ch <- prometheus.MustNewConstMetric(
		numProcesses,
		prometheus.GaugeValue,
		processCount,
	)
}

func main() {
	// Since we are dealing with custom Collector implementations, it might
	// be a good idea to try it out with a pedantic registry.
	reg := prometheus.NewPedanticRegistry()

	pm := processManager{}

	// Construct cluster managers. In real code, we would assign them to
	// variables to then do something with them.
	prometheus.WrapRegistererWith(prometheus.Labels{}, reg).MustRegister(pm)

	// Add the standard process and Go metrics to the custom registry.
	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	http.Handle("/pod-metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8082", nil))
}
