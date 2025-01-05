package describe

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var DescribeJobsCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "opengovernance",
	Subsystem: "scheduler",
	Name:      "describe_jobs_total",
	Help:      "Count of describe jobs",
}, []string{"status"})

var DescribeSourceJobsCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "opengovernance",
	Subsystem: "scheduler",
	Name:      "describe_source_jobs_total",
	Help:      "Count of describe source jobs",
}, []string{"status"})

var DescribeResourceJobsCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "opengovernance",
	Subsystem: "scheduler",
	Name:      "describe_resource_jobs_total",
	Help:      "Count of describe resource jobs",
}, []string{"status", "spot"})

var CleanupJobCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "opengovernance",
	Subsystem: "scheduler",
	Name:      "cleanup_jobs_total",
	Help:      "Count of cleanup jobs",
}, []string{"status"})

var ResourcesDescribedCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "opengovernance",
	Subsystem: "scheduler",
	Name:      "resources_described_total",
	Help:      "Count of resources described",
}, []string{"provider", "status"})

var ResourceBatchProcessLatency = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "opengovernance",
		Subsystem:  "scheduler",
		Name:       "resource_batch_process_millisecond",
		Help:       "Total resource batch process latency.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"provider"})

var ResourceProcessLatency = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "opengovernance",
		Subsystem:  "scheduler",
		Name:       "resource_process_millisecond",
		Help:       "Total resource process latency.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"provider"})

var ResultsDeliveredCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "opengovernance",
	Subsystem: "scheduler",
	Name:      "results_delivered_total",
	Help:      "Count of results delivered",
}, []string{"provider"})

var ResultsProcessedCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "opengovernance",
	Subsystem: "scheduler",
	Name:      "results_processed_total",
	Help:      "Count of results processed",
}, []string{"provider", "status"})

var OldResourcesDeletedCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "opengovernance",
	Subsystem: "scheduler",
	Name:      "old_resources_deleted_total",
	Help:      "Count of old resources deleted",
}, []string{"provider"})

var StreamFailureCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "opengovernance",
	Subsystem: "scheduler",
	Name:      "stream_failure_total",
	Help:      "Count of failures in streams",
}, []string{"provider"})
