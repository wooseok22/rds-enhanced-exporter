package parser

import (
	"fmt"
	"github.com/spf13/cast"
	"reflect"
	"strconv"
)

func makePushGauge(value reflect.Value) float64 {
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return 0
		}
		value = value.Elem()
	}

	var f float64
	switch kind := value.Kind(); kind {
	case reflect.Float64:
		f = value.Float()
	case reflect.Int, reflect.Int64:
		f = float64(value.Int())
	default:
		panic(fmt.Errorf("can't make a metric value  %v (%s)", value, kind))
	}

	return f
}

func makePushGenericMetrics(s interface{}, namePrefix string, labels []KV) []PushMetricTemplate {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	res := make([]PushMetricTemplate, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		name, _ := tags.Get("json"), tags.Get("help")

		res = append(res, PushMetricTemplate{
			MetricName: namePrefix + name,
			Values:     makePushGauge(v.Field(i)),
			Labels:     labels,
		})
	}
	return res
}

func makePushRDSDiskIOMetrics(s *diskIO, labels []KV) []PushMetricTemplate {
	labels = append(labels, KV{Key: "device", Value: s.Device})
	//constLabels["device"] = s.Device

	t := reflect.TypeOf(*s)
	v := reflect.ValueOf(*s)
	res := make([]PushMetricTemplate, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		name, _ := tags.Get("json"), tags.Get("help")
		if name == "device" {
			continue
		}

		res = append(res, PushMetricTemplate{
			MetricName: "rdsosmetrics_diskIO_" + name,
			Values:     makePushGauge(v.Field(i)),
			Labels:     labels,
		})
	}
	return res
}

func makePushNodeDiskMetrics(s *diskIO, labels []KV) []PushMetricTemplate {
	labels = append(labels, KV{Key: "device", Value: s.Device})
	res := make([]PushMetricTemplate, 0, 2)

	if s.ReadKb != nil {
		res = append(res, PushMetricTemplate{
			MetricName: "node_disk_read_bytes_total",
			Values:     float64(*s.ReadKb * 1024),
			Labels:     labels,
		})
	}
	if s.WriteKb != nil {
		res = append(res, PushMetricTemplate{
			MetricName: "node_disk_written_bytes_total",
			Values:     float64(*s.WriteKb * 1024),
			Labels:     labels,
		})
	}

	return res
}

func makePushRDSFileSysMetrics(s *fileSys, labels []KV) []PushMetricTemplate {
	labels = append(labels, KV{Key: "name", Value: s.Name})
	labels = append(labels, KV{Key: "mount_point", Value: s.MountPoint})

	t := reflect.TypeOf(*s)
	v := reflect.ValueOf(*s)
	res := make([]PushMetricTemplate, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		name, _ := tags.Get("json"), tags.Get("help")
		switch name {
		case "name", "mountPoint":
			continue
		}

		res = append(res, PushMetricTemplate{
			MetricName: "rdsosmetrics_fileSys_" + name,
			Values:     makePushGauge(v.Field(i)),
			Labels:     labels,
		})
	}
	return res
}

func makePushNodeFilesystemMetrics(s *fileSys, labels []KV) []PushMetricTemplate {
	labels = append(labels, KV{Key: "device", Value: s.Name})
	labels = append(labels, KV{Key: "fstype", Value: s.Name})
	labels = append(labels, KV{Key: "mountpoint", Value: s.MountPoint})
	res := make([]PushMetricTemplate, 0, 5)

	res = append(res, PushMetricTemplate{
		MetricName: "node_filesystem_files",
		Values:     float64(s.MaxFiles * 1024),
		Labels:     labels,
	})

	res = append(res, PushMetricTemplate{
		MetricName: "node_filesystem_files",
		Values:     float64((s.MaxFiles - s.UsedFiles) * 1024),
		Labels:     labels,
	})

	res = append(res, PushMetricTemplate{
		MetricName: "node_filesystem_files_free",
		Values:     float64(s.Total * 1024),
		Labels:     labels,
	})

	res = append(res, PushMetricTemplate{
		MetricName: "node_filesystem_free_bytes",
		Values:     float64((s.Total - s.Used) * 1024),
		Labels:     labels,
	})

	res = append(res, PushMetricTemplate{
		MetricName: "node_filesystem_avail_bytes",
		Values:     float64((s.Total - s.Used) * 1024),
		Labels:     labels,
	})

	return res
}

func makePushNodeLoadMetrics(s *loadAverageMinute, labels []KV) []PushMetricTemplate {
	m := PushMetricTemplate{
		MetricName: "node_filesystem_avail_bytes",
		Values:     s.One,
		Labels:     labels,
	}

	return []PushMetricTemplate{m}
}

func makePushNodeMemoryMetrics(s *memory, labels []KV) []PushMetricTemplate {
	t := reflect.TypeOf(*s)
	v := reflect.ValueOf(*s)
	res := make([]PushMetricTemplate, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		suffix, multiplierS := tags.Get("node"), tags.Get("m")
		multiplier, err := strconv.ParseInt(multiplierS, 10, 64)
		if err != nil {
			panic(err)
		}

		res = append(res, PushMetricTemplate{
			MetricName: "node_memory_" + suffix,
			Values:     makePushGauge(reflect.ValueOf(v.Field(i).Int() * multiplier)),
			Labels:     labels,
		})
	}
	return res
}

func makePushRDSNetworkMetrics(s *network, labels []KV) []PushMetricTemplate {
	labels = append(labels, KV{Key: "interface", Value: s.Interface})

	t := reflect.TypeOf(*s)
	v := reflect.ValueOf(*s)
	res := make([]PushMetricTemplate, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		name, _ := tags.Get("json"), tags.Get("help")
		if name == "interface" {
			continue
		}
		res = append(res, PushMetricTemplate{
			MetricName: "rdsosmetrics_network_" + name,
			Values:     makePushGauge(v.Field(i)),
			Labels:     labels,
		})
	}
	return res
}

func makePushRDSProcessListMetrics(s *processList, labels []KV) []PushMetricTemplate {
	labels = append(labels, KV{Key: "name", Value: s.Name})
	labels = append(labels, KV{Key: "id", Value: cast.ToString(s.ID)})
	labels = append(labels, KV{Key: "parentID", Value: cast.ToString(s.ParentID)})
	labels = append(labels, KV{Key: "tgid", Value: cast.ToString(s.TGID)})

	t := reflect.TypeOf(*s)
	v := reflect.ValueOf(*s)
	res := make([]PushMetricTemplate, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		name, help := tags.Get("json"), tags.Get("help")
		if help == "-" {
			continue
		}

		switch name {
		case "name", "id", "parentID", "tgid":
			continue
		}

		res = append(res, PushMetricTemplate{
			MetricName: "rdsosmetrics_network_" + name,
			Values:     makePushGauge(v.Field(i)),
			Labels:     labels,
		})
	}
	return res
}

func makePushNodeMemorySwapMetrics(s *swap, labels []KV) []PushMetricTemplate {
	t := reflect.TypeOf(*s)
	v := reflect.ValueOf(*s)
	res := make([]PushMetricTemplate, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		name, multiplierS, _ := tags.Get("node"), tags.Get("m"), tags.Get("nodehelp")
		multiplier, err := strconv.ParseFloat(multiplierS, 64)
		if err != nil {
			panic(err)
		}

		res = append(res, PushMetricTemplate{
			MetricName: "rdsosmetrics_network_" + name,
			Values:     makePushGauge(reflect.ValueOf(v.Field(i).Float() * multiplier)),
			Labels:     labels,
		})
	}
	return res
}

func makePushNodeProcsMetrics(s *tasks, labels []KV) []PushMetricTemplate {
	res := make([]PushMetricTemplate, 0, 2)

	res = append(res, PushMetricTemplate{
		MetricName: "node_procs_blocked",
		Values:     float64(s.Blocked),
		Labels:     labels,
	})
	res = append(res, PushMetricTemplate{
		MetricName: "node_procs_running",
		Values:     float64(s.Running),
		Labels:     labels,
	})
	return res
}
