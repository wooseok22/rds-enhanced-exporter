package cmd

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"rds-enhanced-exporter/config"
	"rds-enhanced-exporter/exporter"
	"rds-enhanced-exporter/log"
)

var rootCMD = &cobra.Command{
	Use: "rds-enhanced-exporter",
	Run: func(cmd *cobra.Command, args []string) {
		confPath, err := cmd.Flags().GetString("conf")
		if err != nil || confPath == "" {
			fmt.Printf("Failed to read configuration file : %v", err)
			os.Exit(1)
		}
		prefix, _ := cmd.Flags().GetString("log")
		config.Initialize(confPath)
		log.SetLogDir(prefix)
		log.Initialize()

		if config.GetConfig().Global.ScrapMethod == "push" {
			exporter.StartPushExporter()
			port := ":" + config.GetConfig().Global.Port
			log.INFO().Infof("[START] started~! (%v)", port)
			http.ListenAndServe(port, nil)
			if err != nil {
				log.ERROR().Fatal("Failed to start server on :", err)
			}
		} else if config.GetConfig().Global.ScrapMethod == "pull" {
			registry := prometheus.NewRegistry()
			registry.MustRegister(exporter.StartPullExporter())

			// default routing
			http.Handle("/", promhttp.HandlerFor(registry, promhttp.HandlerOpts{
				ErrorHandling: promhttp.ContinueOnError,
			}))

			// dummy alias
			http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{
				ErrorHandling: promhttp.ContinueOnError,
			}))

			port := ":" + config.GetConfig().Global.Port
			log.INFO().Infof("[START] started~! (%v)", port)
			http.ListenAndServe(port, nil)
			if err != nil {
				log.ERROR().Fatal("Failed to start server on :", err)
			}
		}
	},
}

func init() {
	rootCMD.Flags().StringP("log", "l", "./", "Set the directory to store all log files")
	rootCMD.Flags().StringP("conf", "c", "./rds-enhanced-exporter.conf", "Set configuration file")
	rootCMD.Flags().StringP("debug", "d", "", "Execute Mode")
}

func Execute() {
	_ = rootCMD.Execute()
}
