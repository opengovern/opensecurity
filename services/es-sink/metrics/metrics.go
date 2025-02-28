package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var EsSinkDocsNumAdded = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "opengovernance",
	Subsystem: "es_sink",
	Name:      "docs_num_added",
	Help:      "Number of documents added to es sink",
})

var EsSinkDocsNumFlushed = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "opengovernance",
	Subsystem: "es_sink",
	Name:      "docs_num_flushed",
	Help:      "Number of documents flushed from es sink",
})

var EsSinkDocsNumFailed = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "opengovernance",
	Subsystem: "es_sink",
	Name:      "docs_num_failed",
	Help:      "Number of documents failed in es sink",
})

var EsSinkDocsNumIndexed = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "opengovernance",
	Subsystem: "es_sink",
	Name:      "docs_num_indexed",
	Help:      "Number of documents indexed in es sink",
})

var EsSinkDocsNumCreated = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "opengovernance",
	Subsystem: "es_sink",
	Name:      "docs_num_created",
	Help:      "Number of documents created in es sink",
})

var EsSinkDocsNumUpdated = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "opengovernance",
	Subsystem: "es_sink",
	Name:      "docs_num_updated",
	Help:      "Number of documents updated in es sink",
})

var EsSinkDocsNumDeleted = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "opengovernance",
	Subsystem: "es_sink",
	Name:      "docs_num_deleted",
	Help:      "Number of documents deleted in es sink",
})

var EsSinkDocsNumRequests = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "opengovernance",
	Subsystem: "es_sink",
	Name:      "docs_num_requests",
	Help:      "Number of requests to es sink",
})
