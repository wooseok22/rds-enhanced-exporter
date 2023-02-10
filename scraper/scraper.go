package scraper

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/prometheus/client_golang/prometheus"
	"rds-enhanced-exporter/config"
	"rds-enhanced-exporter/log"
	"rds-enhanced-exporter/parser"
	"time"
)

func (s *Scraper) PullStart(done chan bool, svc *cloudwatchlogs.CloudWatchLogs, ch chan<- map[string][]prometheus.Metric) {
	ticker := time.NewTicker(config.GetConfig().Global.TInterval * time.Second)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.ERROR().Errorf("scraper recover : %v", r)
				s.PullStart(done, svc, ch)
				done <- false
			}
		}()
		for {
			select {
			case <-ticker.C:
				metrics := s.Scrap(svc)
				ch <- metrics
			}
		}
	}()

	if <-done {
	}
}

func (s *Scraper) PushStart(done chan bool, svc *cloudwatchlogs.CloudWatchLogs, ch chan<- map[string][]prometheus.Metric) {
	ticker := time.NewTicker(config.GetConfig().Global.TInterval * time.Second)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.ERROR().Errorf("scraper recover : %v", r)
				s.PushStart(done, svc, ch)
				done <- false
			}
		}()
		for {
			select {
			case <-ticker.C:
				_ = s.Scrap(svc)
			}
		}
	}()

	if <-done {
	}
}

func (s *Scraper) Scrap(svc *cloudwatchlogs.CloudWatchLogs) map[string][]prometheus.Metric {
	var osMetrics []parser.OSMetrics
	var osMetric *parser.OSMetrics
	var resultMetrics map[string][]prometheus.Metric

	queryOutput, err := svc.StartQuery(GetQueryInput())
	if err != nil {
		log.ERROR().Errorf("[Scraper] StartQuery err : %v", err)
	}

	time.Sleep(3 * time.Second)

	queryResult, err := svc.GetQueryResults(GetQueryResult(*queryOutput.QueryId))
	if err != nil {
		log.ERROR().Errorf("[Scraper] GetQueryResults err : %v", err)
	}

	for _, v := range queryResult.Results {
		// one osMetric = one instance metrics
		osMetric = nil
		d := json.NewDecoder(bytes.NewReader([]byte(*v[0].Value)))
		if err := d.Decode(&osMetric); err != nil {
			log.INFO().Errorf("[Scraper] JSON decode err : %v", err)
		}
		osMetrics = append(osMetrics, *osMetric)
	}

	if config.GetConfig().Global.ScrapMethod == "push" {
		PushMetric(parser.MakePushPrometheusMetrics(osMetrics))
	} else {
		resultMetrics = parser.MakePullPrometheusMetrics(osMetrics)
	}

	return resultMetrics
}
