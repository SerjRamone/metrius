// Package main is agent main package
package main

import (
	"log"
	"time"

	collect "github.com/SerjRamone/metrius/internal/collector"
	"github.com/SerjRamone/metrius/internal/config"
	"github.com/SerjRamone/metrius/internal/sender"
	"github.com/SerjRamone/metrius/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	conf, err := config.NewAgent()
	if err != nil {
		log.Fatal("config parse error: ", err)
	}

	if err := logger.Init("info"); err != nil {
		log.Fatal("can't init logger")
	}

	logger.Info("loaded config", zap.Object("config", &conf))

	reportedAt := time.Now()
	polledAt := time.Now()

	sender := sender.NewMetricsSender(conf.ServerAddress)
	collector := collect.New()

	for {
		// collect metrics
		seconds := int((time.Since(polledAt)).Seconds())
		if seconds >= conf.PollInterval {
			collector.Collect()
			polledAt = time.Now()
		}

		// send metrics
		seconds = int((time.Since(reportedAt)).Seconds())
		if seconds >= conf.ReportInterval {
			if collections := collector.Export(); len(collections) > 0 {
				err := sender.Send(collections)
				if err != nil {
					logger.Error("sender error", zap.Error(err))
				}
			}

			reportedAt = time.Now()
		}

		time.Sleep(500 * time.Millisecond)
	}
}
