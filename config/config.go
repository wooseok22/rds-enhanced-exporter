package config

import (
	"fmt"
	"gopkg.in/gcfg.v1"
	"os"
)

var (
	info ConfInfo
)

func Initialize(confPath string) {
	err := info.ReadConfig(confPath)

	if err != nil {
		fmt.Printf("Failed to parse conf data : %s", err)
		os.Exit(1)
	}
}

func GetConfig() ConfInfo {
	return info
}
func (confInfo *ConfInfo) ReadConfig(confPath string) error {
	var cfg ConfInfo

	err := gcfg.ReadFileInto(&cfg, confPath)
	if err != nil {
		return err
	}

	confInfo.Global.AwsArn = cfg.Global.AwsArn
	confInfo.Global.Region = cfg.Global.Region
	confInfo.Global.Port = cfg.Global.Port
	confInfo.Global.SInterval = cfg.Global.SInterval
	confInfo.Global.TInterval = cfg.Global.TInterval
	confInfo.Global.Log = cfg.Global.Log
	confInfo.Global.ScrapMethod = cfg.Global.ScrapMethod
	confInfo.Global.TSDBUrl = cfg.Global.TSDBUrl
	confInfo.Labels.Target = cfg.Labels.Target
	confInfo.Labels.Kv = cfg.Labels.Kv

	return nil
}
