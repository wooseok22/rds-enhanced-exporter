# RDS Enhanced Exporter 
## Information

The more RDS instances managed by one aws account, the more likely "CloudWatchLogs filterLogEvents" will cause AWS API throttling. (percona rds-exporter's method)

The program uses the "CloudWatchLogs queryStart" API to make only one call in a single collection cycle to avoid AWS API-throttling.

This program also allows you to collect metrics for all instances with only one exporter agent in one AWS account.

## Architecture
### pull
<img width="900" alt="pull_architecture" src="https://user-images.githubusercontent.com/21119703/218067081-c7fda7e9-0af7-46de-87df-4a2bc85096a3.png">

### push
<img width="900" alt="push_architecture" src="https://user-images.githubusercontent.com/21119703/218067099-4e102d58-167a-4f01-9c22-e9a4d2f5797b.png">


## Runtime Option
| Parameter   | Description                 |
|-------------|-----------------------------|
| `--log -l`  | set log file path           |
| `--conf -c` | set configuration file path |


## Configuration file information
* `o is required information`
* `x is not required information`
### global
| Parameter                     | Description                                                                  |
|-------------------------------|------------------------------------------------------------------------------|
| `awsarn[string]`(o)           | ARN of AWS accounts to collect                                               |
| `region[string]`(o)           | Region of AWS accounts to collect                                            |
| `port[string]`(o)             | Port of rds enhanced exporter                                                |
| `sinterval[time.duration]`(o) | Scrap interval time (data from a few minutes ago) (recommended settings:3)   |
| `tinterval[time.duration]`(o) | Ticker interval time (how many seconds to collect) (recommended settings:15) |
| `log[string]`(x)              | Program log file name                                                        |
| `sentrydsn[string]`(x)        | Set this information to receive error logs via sentry                        |
| `scrapmethod`(o)              | Scrap Method (push or pull / default : pull)                                 |
| `tsdburl`(x(*))               | timeSeries database endpoint to use when pushing                             |


### labels
| Parameter             | Description                                                                                                                |
|-----------------------|----------------------------------------------------------------------------------------------------------------------------|
| `target`(o)           | Value to enter instead of "instance" label                                                                                 |
| `kv[string array]`(x) | This means "Key/Value". You can set metrics labels through this option. The delimiter between key and value is ":(colon)". |
* labels config exmaple
  * target : host_name / kv[0] : region:ap-northeast-2 / kv[1] : account:db-prod
    * metric : rdsosmetrics_timestamp{host_name="test-db01", region="ap-northeast-2", account="db-prod", job="rds-enhanced"}
  * target : instance / kv[0] : region:ap-northeast-2 / kv[1] : account_name:db-prod
    * metric : rdsosmetrics_timestamp{instance="test-db01", region="ap-northeast-2", account_name="db-prod", job="rds-enhanced"}



## How to build (example)
```
$ git clone thisProject

$ cd ./thisProject

$ go build

$ ./rds-enhanced-exporter -l ./program_logs -c ./rds-enhanced-exporter.conf
```


## Limitation
1. This program can support for only AWS RDS enhanced monitoring.
2. Because this program uses the Go AWS SDK, it requires iam role and permissions to run the CloudWatch Logs API.
3. Until now, it only supports authentication methods through iam role ARN. (to be updated later)

