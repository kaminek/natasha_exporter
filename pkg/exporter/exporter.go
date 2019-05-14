package exporter

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"github.com/kaminek/natasha-cli/pkg/handlers"
	"github.com/kaminek/natasha-cli/pkg/headers"
	"github.com/kaminek/natasha_exporter/pkg/config"
	"github.com/kaminek/natasha_exporter/pkg/info"
	"github.com/kaminek/natasha_exporter/pkg/metrics"
	"github.com/kaminek/natasha_exporter/pkg/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	// namespace defines the namespace for the exporter.
	namespace = "natasha"
)

// NatashaCollector collects metrics, mostly runtime, about this exporter in general.
type NatashaCollector struct {
	version   string
	revision  string
	buildDate string
	goVersion string
	startTime time.Time

	BuildInfo *prometheus.Desc
}

// NewNatashaCollectorMetrics returns a collector about the collector itself.
func NewNatashaCollectorMetrics() NatashaCollector {
	return NatashaCollector{
		version:   info.Version,
		revision:  info.Revision,
		buildDate: info.BuildDate,
		goVersion: info.GoVersion,
		startTime: info.StartTime,

		BuildInfo: prometheus.NewDesc(
			"natasha_build_info",
			"Exporter built information.",
			[]string{"version", "revision", "builddate", "goversion"},
			nil,
		),
	}
}

// Exporter main struct
type Exporter struct {
	URI   string
	mutex sync.RWMutex

	totalScrapes            prometheus.Counter
	NatashaStatus           *prometheus.Desc
	NatashaCollectorMetrics NatashaCollector
	AppStatsMetrics         map[string]*prometheus.Desc
	DpdkStatsMetrics        map[string]*prometheus.Desc
}

// NewExporter returns an initialized Exporter.
// func NewExporter(NatashaMetrics map[int]*prometheus.Desc, timeout time.Duration) (*Exporter, error) {
func NewExporter(timeout time.Duration) (*Exporter, error) {
	return &Exporter{
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_total_scrapes",
			Help:      "Current total Natasha scrapes.",
		}),
		NatashaStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"The natasha server status.",
			nil,
			nil),
		NatashaCollectorMetrics: NewNatashaCollectorMetrics(),
		AppStatsMetrics:         metrics.GetAppStatsMetrics(),
		DpdkStatsMetrics:        metrics.GetDpdkStatsMetrics(),
	}, nil
}

func scrape(ch chan<- prometheus.Metric) (up float64) {

	conn, err := server.NatashaServerDial()
	if err != nil {
		return 0
	}

	reply := headers.NatashaCmdReply{}

	// get version
	err = handlers.SendCmd(conn, headers.NatashaCmdVersion, &reply)
	if err != nil {
		log.Fatal("version: ", err)
		return 0
	}

	recvBuf := make([]byte, reply.DataSize)
	_, err = conn.Read(recvBuf)
	if err != nil {
		log.Fatal("Connection error", err)
		return 0
	}
	// infos.version = recvBuf

	err = handlers.SendCmd(conn, headers.NatashaCmdDpdkStats, &reply)
	if err != nil {
		log.Fatal("dpdk-stats: ", err)
		return 0
	}

	dpdkStats := headers.NatashaEthStats{}
	ports := int(reply.DataSize) / int(unsafe.Sizeof(dpdkStats))
	recvBuf = make([]byte, unsafe.Sizeof(dpdkStats))

	dpdkMetrics := metrics.GetDpdkStatsMetrics()

	for port := 0; port < ports; port++ {
		_, err = conn.Read(recvBuf)
		if err != nil {
			log.Fatal("Failed to read data", err)
			return 0
		}

		// Write byte stream to struct
		r := bytes.NewReader(recvBuf)

		err = binary.Read(r, binary.BigEndian, &dpdkStats)
		if err != nil {
			log.Fatal("Write to data structure error: ", err)
			return 0
		}

		ch <- prometheus.MustNewConstMetric(
			dpdkMetrics["Ipackets"],
			prometheus.GaugeValue,
			float64(dpdkStats.Ipackets),
			strconv.Itoa(port),
		)
		ch <- prometheus.MustNewConstMetric(
			dpdkMetrics["Opackets"],
			prometheus.GaugeValue,
			float64(dpdkStats.Opackets),
			strconv.Itoa(port),
		)
		ch <- prometheus.MustNewConstMetric(
			dpdkMetrics["Ibytes"],
			prometheus.GaugeValue,
			float64(dpdkStats.Ibytes),
			strconv.Itoa(port),
		)
		ch <- prometheus.MustNewConstMetric(
			dpdkMetrics["Obytes"],
			prometheus.GaugeValue,
			float64(dpdkStats.Obytes),
			strconv.Itoa(port),
		)
	}

	err = handlers.SendCmd(conn, headers.NatashaCmdAppStats, &reply)
	if err != nil {
		log.Fatal("app-stats: ", err)
		return 0
	}

	var coreID uint8
	appStats := headers.NatashaAppStats{}
	cores := int(reply.DataSize) /
		int(unsafe.Sizeof(appStats)+unsafe.Sizeof(coreID))
	appMetrics := metrics.GetAppStatsMetrics()

	for core := 0; core < cores; core++ {
		// Get coreid
		recvBuf = make([]byte, unsafe.Sizeof(coreID))
		_, err := conn.Read(recvBuf)
		if err != nil {
			log.Fatal("Failed to read data", err)
			return 0
		}
		// it's a uint8 same as one byte
		coreID = recvBuf[0]

		// Get app stats for that core
		recvBuf = make([]byte, unsafe.Sizeof(appStats))
		_, err = conn.Read(recvBuf)
		if err != nil {
			log.Fatal("Failed to read data", err)
			return 0
		}

		r := bytes.NewReader(recvBuf)
		err = binary.Read(r, binary.BigEndian, &appStats)
		if err != nil {
			log.Fatal("Write to data structure error: ", err)
			return 0
		}

		ch <- prometheus.MustNewConstMetric(
			appMetrics["DropNoRule"],
			prometheus.GaugeValue,
			float64(appStats.DropNoRule),
			strconv.Itoa(int(coreID)),
		)
		ch <- prometheus.MustNewConstMetric(
			appMetrics["DropNatCondition"],
			prometheus.GaugeValue,
			float64(appStats.DropNatCondition),
			strconv.Itoa(int(coreID)),
		)
		ch <- prometheus.MustNewConstMetric(
			appMetrics["DropBadL3Cksum"],
			prometheus.GaugeValue,
			float64(appStats.DropBadL3Cksum),
			strconv.Itoa(int(coreID)),
		)
		ch <- prometheus.MustNewConstMetric(
			appMetrics["RxBadL4Cksum"],
			prometheus.GaugeValue,
			float64(appStats.RxBadL4Cksum),
			strconv.Itoa(int(coreID)),
		)
		ch <- prometheus.MustNewConstMetric(
			appMetrics["DropUnknownIcmp"],
			prometheus.GaugeValue,
			float64(appStats.DropUnknownIcmp),
			strconv.Itoa(int(coreID)),
		)
		ch <- prometheus.MustNewConstMetric(
			appMetrics["DropUnhandledEthertype"],
			prometheus.GaugeValue,
			float64(appStats.DropUnhandledEthertype),
			strconv.Itoa(int(coreID)),
		)
		ch <- prometheus.MustNewConstMetric(
			appMetrics["DropTxNotsent"],
			prometheus.GaugeValue,
			float64(appStats.DropTxNotsent),
			strconv.Itoa(int(coreID)),
		)
	}

	conn.Close()

	return 1
}

// Describe describes all the metrics ever exported by the HAProxy exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.totalScrapes.Desc()
	ch <- e.NatashaStatus
	ch <- e.NatashaCollectorMetrics.BuildInfo

	for _, m := range e.AppStatsMetrics {
		ch <- m
	}

	for _, m := range e.DpdkStatsMetrics {
		ch <- m
	}
}

// Collect fetches the stats from configured HAProxy location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	e.totalScrapes.Inc()
	scrape(ch)
	ch <- e.totalScrapes
	ch <- prometheus.MustNewConstMetric(
		e.NatashaCollectorMetrics.BuildInfo,
		prometheus.GaugeValue,
		1.0,
		e.NatashaCollectorMetrics.version,
		e.NatashaCollectorMetrics.revision,
		e.NatashaCollectorMetrics.buildDate,
		e.NatashaCollectorMetrics.goVersion,
	)
}

// New Exporter
func New(cfg *config.Config) {
	exporter, err := NewExporter(cfg.Target.Timeout)
	if err != nil {
		log.Fatal(err)
	}
	prometheus.MustRegister(exporter)

}
