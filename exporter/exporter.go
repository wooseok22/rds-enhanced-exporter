package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"rds-enhanced-exporter/common"
	"rds-enhanced-exporter/log"
	"rds-enhanced-exporter/scraper"
	"time"
)

type exporter struct {
	Metrics map[string][]prometheus.Metric
}

func StartPullExporter() *exporter {
	s := new(scraper.Scraper)
	c := &exporter{Metrics: make(map[string][]prometheus.Metric)}

	ch := make(chan map[string][]prometheus.Metric)
	svc := common.GetAwsCloudWatchLogsSession()

	m := s.Scrap(svc)
	c.SetMetrics(m)
	go func() {
		for m := range ch {
			c.SetMetrics(m)
		}
	}()
	go s.PullStart(make(chan bool), svc, ch)

	return c
}

func StartPushExporter() {
	s := new(scraper.Scraper)
	ch := make(chan map[string][]prometheus.Metric)
	svc := common.GetAwsCloudWatchLogsSession()

	s.PushStart(make(chan bool), svc, ch)
	go s.PushStart(make(chan bool), svc, ch)
}

func (c *exporter) SetMetrics(m map[string][]prometheus.Metric) {
	for id, metric := range m {
		c.Metrics[id] = metric
	}
}

func (c *exporter) Collect(ch chan<- prometheus.Metric) {
	count := 0
	for _, metrics := range c.Metrics {
		for _, metric := range metrics {
			ch <- metric
		}
		count++
	}
	log.INFO().Infof("[SCRAP:%v] uploaded metric count : %v", time.Now().Format("2006-01-02 15:04:05"), count)
}

func (c *exporter) Describe(ch chan<- *prometheus.Desc) {
	// just for prometheus interface
}
