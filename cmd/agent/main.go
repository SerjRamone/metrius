// Package main is agent main package
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"

	collect "github.com/SerjRamone/metrius/internal/collector"
	"github.com/SerjRamone/metrius/internal/config"
	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/internal/sender"
	"github.com/SerjRamone/metrius/pkg/logger"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printTags()

	conf, err := config.NewAgent()
	if err != nil {
		log.Fatal("config parse error: ", err)
	}

	if err = logger.Init("info"); err != nil {
		log.Fatal("can't init logger")
	}

	logger.Info("loaded config", zap.Object("config", &conf))

	var pubKey []byte
	// check path and read key
	if conf.CryptoKey != "" {
		pubKey, err = os.ReadFile(conf.CryptoKey)
		if err != nil {
			logger.Error("reading keyfile error", zap.Error(err))
			return
		}
	}
	sender := sender.NewMetricsSender(conf.ServerAddress, conf.HashKey, pubKey)
	collector := collect.New()

	// closing channel
	doneCh := make(chan struct{})

	// catch signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// chan for jobs for senders
	jobCh := make(chan []metrics.Collection)

	// run workers
	for w := 1; w <= conf.RateLimit; w++ {
		go sender.Worker(doneCh, jobCh)
	}

	// collect metrics
	go func() {
		ticker := time.NewTicker(time.Duration(conf.PollInterval) * time.Second)
		for {
			select {
			case <-ticker.C:
				collector.Collect()

			case <-doneCh:
				logger.Info("collector recived done signal")
				return
			}
		}
	}()

	// additional metrics
	go func() {
		ticker := time.NewTicker(time.Duration(conf.ReportInterval) * time.Second)
		for {
			select {
			case <-ticker.C:
				logger.Info("collect additional metrics")
				v, err := mem.VirtualMemory()
				if err != nil {
					logger.Error("getting memory metrics error", zap.Error(err))
					continue
				}
				cpu, err := cpu.Percent(0, true)
				if err != nil {
					logger.Error("getting cpu metrics error", zap.Error(err))
					continue
				}
				c := metrics.Collection{
					metrics.CollectionItem{Name: "TotalMemory", Type: "gauge", Value: float64(v.Total)},
					metrics.CollectionItem{Name: "FreeMemory", Type: "gauge", Value: float64(v.Free)},
					metrics.CollectionItem{Name: "CPUutilization1", Type: "gauge", Value: cpu[0]},
				}
				collections := []metrics.Collection{c}
				jobCh <- collections

			case <-doneCh:
				logger.Info("additional metrics recived done signal")
				return
			}
		}
	}()

	// put jobs
	go func() {
		defer close(jobCh)
		ticker := time.NewTicker(time.Duration(conf.ReportInterval) * time.Second)
		for {
			select {
			case <-ticker.C:
				if collections := collector.Export(); len(collections) > 0 {
					jobCh <- collections
				}

			case <-doneCh:
				logger.Info("sender recived done signal")
				if collections := collector.Export(); len(collections) > 0 {
					jobCh <- collections
				}
				return
			}
		}
	}()

	<-sigCh
	close(doneCh)
	// waiting gorutines to stop
	// maybe need to use WaitGroup
	time.Sleep(1 * time.Second)
	logger.Info("shutting down")
}

func printTags() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
