package scraper

import "C"
import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"io"
	"net/http"
	"rds-enhanced-exporter/config"
	"rds-enhanced-exporter/log"
	"rds-enhanced-exporter/parser"
	"strings"
	"time"
)

func GetQueryInput() *cloudwatchlogs.StartQueryInput {
	now := time.Now()
	endTime := now.UnixNano() / 1000000
	startTime := now.Add(-(config.GetConfig().Global.SInterval * time.Minute)).UnixNano() / 1000000
	input := &cloudwatchlogs.StartQueryInput{
		StartTime:    aws.Int64(startTime),
		EndTime:      aws.Int64(endTime),
		LogGroupName: aws.String("RDSOSMetrics"),
		QueryString:  aws.String(`fields @message`),
	}

	return input
}

func GetQueryResult(queryId string) *cloudwatchlogs.GetQueryResultsInput {
	return &cloudwatchlogs.GetQueryResultsInput{QueryId: aws.String(queryId)}
}

func PushMetric(pushData []parser.PushMetricTemplate) {
	for _, v1 := range pushData {
		client := &http.Client{}
		labels := fmt.Sprintf("{")
		for i, v2 := range v1.Labels {
			if i == len(v1.Labels)-1 {
				labels += v2.Key + `="` + v2.Value + `"}`
				break
			}
			labels += v2.Key + `="` + v2.Value + `",`
		}

		var str string
		switch v := v1.Values.(type) {
		default:
			fmt.Printf("unexpected type %T", v)
		case float64:
			str = fmt.Sprintf("%s%s %v", v1.MetricName, labels, v1.Values)
		case float32:
			str = fmt.Sprintf("%s%s %v", v1.MetricName, labels, v1.Values)
		case int64:
			str = fmt.Sprintf("%s%s %v", v1.MetricName, labels, v1.Values)
		case int32:
			str = fmt.Sprintf("%s%s %v", v1.MetricName, labels, v1.Values)
		case string:
			str = fmt.Sprintf("%s%s %s", v1.MetricName, labels, v1.Values)
		}

		var data = strings.NewReader(str)
		req, err := http.NewRequest("POST", config.GetConfig().Global.TSDBUrl, data)
		if err != nil {
			log.ERROR().Errorf("[PUSH-METRIC] push error %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		if err != nil {
			log.ERROR().Errorf("[PUSH-METRIC] push error %v", err)
		}

		_, err = io.ReadAll(resp.Body)
		if err != nil {
			log.ERROR().Errorf("[PUSH-METRIC] push error %v", err)
		}
	}
}
