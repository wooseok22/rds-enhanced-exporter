package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"rds-enhanced-exporter/config"
)

func GetAwsCloudWatchLogsSession() *cloudwatchlogs.CloudWatchLogs {
	sess := session.Must(session.NewSession())
	sess = session.New(&aws.Config{
		Credentials: stscreds.NewCredentials(sess, config.GetConfig().Global.AwsArn),
		Region:      aws.String(config.GetConfig().Global.Region),
	})
	return cloudwatchlogs.New(sess)
}
