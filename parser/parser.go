package parser

import (
	"github.com/prometheus/client_golang/prometheus"
	"rds-enhanced-exporter/common"
	"rds-enhanced-exporter/config"
)

func MakePullPrometheusMetrics(metrics []OSMetrics) map[string][]prometheus.Metric {
	resultMetrics := make(map[string][]prometheus.Metric)
	instanceList := make(map[string]float64)

	for _, metric := range metrics {
		if (instanceList[metric.InstanceID] > float64(metric.Timestamp.Unix())) || instanceList[metric.InstanceID] != 0 {
			continue
		} else {
			var m []prometheus.Metric
			constLabels := prometheus.Labels{
				config.GetConfig().Labels.Target: metric.InstanceID,
				"job":                            "rds-enhanced",
			}
			for _, kv := range config.GetConfig().Labels.Kv {
				constLabels[common.GetKey(kv)] = common.GetValue(kv)
			}

			m = append(m, prometheus.MustNewConstMetric(
				prometheus.NewDesc("rdsosmetrics_timestamp", "Metrics timestamp (UNIX seconds).", nil, constLabels),
				prometheus.CounterValue,
				float64(metric.Timestamp.Unix())),
			)

			m = append(m, prometheus.MustNewConstMetric(
				prometheus.NewDesc("rdsosmetrics_General_numVCPUs", "The number of virtual CPUs for the DB instance.", nil, constLabels),
				prometheus.GaugeValue,
				float64(metric.NumVCPUs)),
			)

			metrics := makeGenericMetrics(metric.CPUUtilization, "rdsosmetrics_cpuUtilization_", constLabels)

			m = append(m, metrics...)
			metrics = makeNodeCPUMetrics(&metric.CPUUtilization, constLabels)
			m = append(m, metrics...)

			for _, disk := range metric.DiskIO {
				metrics = makeRDSDiskIOMetrics(&disk, constLabels)
				m = append(m, metrics...)
				metrics = makeNodeDiskMetrics(&disk, constLabels)
				m = append(m, metrics...)
			}

			for _, fs := range metric.FileSys {
				metrics = makeRDSFileSysMetrics(&fs, constLabels)
				m = append(m, metrics...)
				metrics = makeNodeFilesystemMetrics(&fs, constLabels)
				m = append(m, metrics...)
			}

			metrics = makeGenericMetrics(metric.LoadAverageMinute, "rdsosmetrics_loadAverageMinute_", constLabels)
			m = append(m, metrics...)
			metrics = makeNodeLoadMetrics(&metric.LoadAverageMinute, constLabels)
			m = append(m, metrics...)

			metrics = makeGenericMetrics(metric.Memory, "rdsosmetrics_memory_", constLabels)
			m = append(m, metrics...)
			metrics = makeNodeMemoryMetrics(&metric.Memory, constLabels)
			m = append(m, metrics...)

			for _, n := range metric.Network {
				metrics = makeRDSNetworkMetrics(&n, constLabels)
				m = append(m, metrics...)
			}

			for _, p := range metric.ProcessList {
				metrics = makeRDSProcessListMetrics(&p, constLabels)
				m = append(m, metrics...)
			}

			metrics = makeGenericMetrics(metric.Swap, "rdsosmetrics_swap_", constLabels)
			m = append(m, metrics...)
			metrics = makeNodeMemorySwapMetrics(&metric.Swap, constLabels)
			m = append(m, metrics...)

			metrics = makeGenericMetrics(metric.Tasks, "rdsosmetrics_tasks_", constLabels)
			m = append(m, metrics...)
			metrics = makeNodeProcsMetrics(&metric.Tasks, constLabels)
			m = append(m, metrics...)
			resultMetrics[metric.InstanceID] = m

		}

	}

	return resultMetrics
}

func MakePushPrometheusMetrics(metrics []OSMetrics) []PushMetricTemplate {
	var pushMetrics []PushMetricTemplate
	instanceList := make(map[string]float64)

	for _, metric := range metrics {
		if (instanceList[metric.InstanceID] > float64(metric.Timestamp.Unix())) || instanceList[metric.InstanceID] != 0 {
			continue
		} else {
			var labels []KV
			labels = append(labels, KV{config.GetConfig().Labels.Target, metric.InstanceID})
			for _, kv := range config.GetConfig().Labels.Kv {
				labels = append(labels, KV{common.GetKey(kv), common.GetValue(kv)})
			}

			pushMetrics = append(pushMetrics, PushMetricTemplate{
				MetricName: "rdsosmetrics_timestamp",
				Values:     float64(metric.Timestamp.Unix()),
				Labels:     labels,
			})

			pushMetrics = append(pushMetrics, PushMetricTemplate{
				MetricName: "rdsosmetrics_General_numVCPUs",
				Values:     float64(metric.NumVCPUs),
				Labels:     labels,
			})

			metrics := makePushGenericMetrics(metric.CPUUtilization, "rdsosmetrics_cpuUtilization_", labels)
			pushMetrics = append(pushMetrics, metrics...)

			for _, disk := range metric.DiskIO {
				metrics = makePushRDSDiskIOMetrics(&disk, labels)
				pushMetrics = append(pushMetrics, metrics...)
				metrics = makePushNodeDiskMetrics(&disk, labels)
				pushMetrics = append(pushMetrics, metrics...)
			}

			for _, fs := range metric.FileSys {
				metrics = makePushRDSFileSysMetrics(&fs, labels)
				pushMetrics = append(pushMetrics, metrics...)
				metrics = makePushNodeFilesystemMetrics(&fs, labels)
				pushMetrics = append(pushMetrics, metrics...)
			}

			metrics = makePushGenericMetrics(metric.LoadAverageMinute, "rdsosmetrics_loadAverageMinute_", labels)
			pushMetrics = append(pushMetrics, metrics...)
			metrics = makePushNodeLoadMetrics(&metric.LoadAverageMinute, labels)
			pushMetrics = append(pushMetrics, metrics...)

			metrics = makePushGenericMetrics(metric.Memory, "rdsosmetrics_memory_", labels)
			pushMetrics = append(pushMetrics, metrics...)
			metrics = makePushNodeMemoryMetrics(&metric.Memory, labels)
			pushMetrics = append(pushMetrics, metrics...)

			for _, n := range metric.Network {
				metrics = makePushRDSNetworkMetrics(&n, labels)
				pushMetrics = append(pushMetrics, metrics...)
			}

			for _, p := range metric.ProcessList {
				metrics = makePushRDSProcessListMetrics(&p, labels)
				pushMetrics = append(pushMetrics, metrics...)
			}

			metrics = makePushGenericMetrics(metric.Swap, "rdsosmetrics_swap_", labels)
			pushMetrics = append(pushMetrics, metrics...)
			metrics = makePushNodeMemorySwapMetrics(&metric.Swap, labels)
			pushMetrics = append(pushMetrics, metrics...)

			metrics = makePushGenericMetrics(metric.Tasks, "rdsosmetrics_tasks_", labels)
			pushMetrics = append(pushMetrics, metrics...)
			metrics = makePushNodeProcsMetrics(&metric.Tasks, labels)
			pushMetrics = append(pushMetrics, metrics...)
		}
	}

	return pushMetrics
}
