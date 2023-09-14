// Package main is agent main package
package main

import (
	"log"
	"time"

	"github.com/SerjRamone/metrius/internal/collector"
	"github.com/SerjRamone/metrius/internal/config"
	"github.com/SerjRamone/metrius/internal/sender"
)

func main() {
	conf := config.Agent{}
	conf.ParseFlags()
	err := conf.ParseEnv()
	if err != nil {
		log.Fatal("config parse error: ", err)
	}

	log.Printf("Loaded agent config: %+v\n", conf)

	start(conf)
}

func start(conf config.Agent) {
	reportedAt := time.Now()
	polledAt := time.Now()

	sender := sender.NewMetricsSender(conf.ServerAddress)
	collector := collector.NewCollector()

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
					log.Println("sender error: ", err)
				}
				log.Println("📨  sended")
			}

			reportedAt = time.Now()
		}

		time.Sleep(500 * time.Millisecond)
	}
}
