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
		"Ipackets": newDpdkStatsMetric("ipackets", "The number of Rx packets.", nil),
		"Opackets": newDpdkStatsMetric("opackets", "The number of Tx packets.", nil),
		"Ibytes":   newDpdkStatsMetric("ibytes", "The number of Rx bytes.", nil),
		"Obytes":   newDpdkStatsMetric("obytes", "The number of Tx bytes.", nil),
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
