package scraper

import (
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"time"
)

type Scraper struct {
	LogGroupName string
	startTime    time.Time
	endTime      time.Time
	svc          *cloudwatchlogs.CloudWatchLogs
}
