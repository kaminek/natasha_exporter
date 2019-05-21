package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// namespace defines the namespace for the exporter.
	namespace = "natasha"
)

var (
	dpdkLabelNames = []string{"portid"}
	appLabelNames  = []string{"coreid"}

	// LastScrapeStatus describe the last scrape status
	LastScrapeStatus = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "last_scrape_status"),
		"Was the last scrape successful.",
		nil,
		nil)
)

type metrics map[string]*prometheus.Desc

func newDpdkStatsMetric(metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "dpdk_stats", metricName), docString, dpdkLabelNames, constLabels)
}

func newAppStatsMetric(metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "app_stats", metricName), docString, appLabelNames, constLabels)
}

// GetDpdkStatsMetrics returns all metrics
func GetDpdkStatsMetrics() map[string]*prometheus.Desc {
	return metrics{
		"Ipackets": newDpdkStatsMetric("ipackets", "Total number of successfully received packets.", nil),
		"Opackets": newDpdkStatsMetric("opackets", "Total number of successfully transmitted packets.", nil),
		"Ibytes":   newDpdkStatsMetric("ibytes", "Total number of successfully received bytes.", nil),
		"Obytes":   newDpdkStatsMetric("obytes", "Total number of successfully transmitted bytes.", nil),
		"Imissed":  newDpdkStatsMetric("imissed", "Total of RX packets dropped by the HW because there are no available buffer (i.e. RX queues are full).", nil),
		"Ierrors":  newDpdkStatsMetric("ierrors", "Total number of erroneous received packets.", nil),
		"Oerrors":  newDpdkStatsMetric("oerrors", "Total number of failed transmitted packets.", nil),
		"RxNombuf": newDpdkStatsMetric("rxnombuf", "Total number of RX mbuf allocation failures.", nil),
	}

}

// GetAppStatsMetrics returns all metrics
func GetAppStatsMetrics() map[string]*prometheus.Desc {
	return metrics{
		"DropNoRule":             newAppStatsMetric("drop_no_rule", "The number of Rx packets dropped due to not NAT rule.", nil),
		"DropNatCondition":       newAppStatsMetric("drop_nat_condition", "The number of drops due to ip range missmatch.", nil),
		"DropBadL3Cksum":         newAppStatsMetric("drop_bad_l3_cksum", "The number of Rx packets dropped due bad l3 checksum.", nil),
		"RxBadL4Cksum":           newAppStatsMetric("rx_bad_l4_cksum", "The number of Rx packets having a bad TCP or UDP checksum.", nil),
		"DropUnknownIcmp":        newAppStatsMetric("drop_unknown_icmp", "The number of Rx packets dropped due unhandled or unknow ICMP type.", nil),
		"DropUnhandledEthertype": newAppStatsMetric("drop_unhandled_ether_type", "The number of Rx packets dropped due to unhandled ether type.", nil),
		"DropTxNotsent":          newAppStatsMetric("drop_tx_not_sent", "The number of failed Tx packets.", nil),
	}

}

// GetCPUUsageMetrics returns real cpu usage metric
func GetCPUUsageMetrics() *prometheus.Desc {
	return newAppStatsMetric("real_cpu_usage", "The real CPU usage.", nil)
}
