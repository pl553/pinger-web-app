package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"pinger/internal/api"
	"pinger/internal/docker"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	probing "github.com/prometheus-community/pro-bing"
	"github.com/spf13/cobra"
)

const containerStatusUP = "UP"
const containerStatusDOWN = "DOWN"

var (
	backendUrl    string
	pingPeriodSec int

	rootCmd = &cobra.Command{
		Use:     "pinger",
		Short:   "pinger service",
		Run:     Run,
		PreRunE: PreRunE,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&backendUrl, "backend_url",
		"http://localhost:8091",
		"Backend URL")
	rootCmd.PersistentFlags().IntVar(&pingPeriodSec, "ping_period_sec",
		1,
		"Perform docker container pings every N seconds")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func PreRunE(cmd *cobra.Command, args []string) error {
	if pingPeriodSec <= 0 {
		return fmt.Errorf("ping_period_sec: must be > 0")
	}

	return nil
}

func Run(cmd *cobra.Command, args []string) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	ticker := time.NewTicker(time.Duration(pingPeriodSec) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ips, err := docker.GetRunningContainerIps()
			if err != nil {
				log.Fatal(err)
			}

			for _, ip := range ips {
				go pingAndReport(ip)
			}
		case <-sigChan:
			return
		}
	}
}

func pingAndReport(ip string) {
	reportDown := func() {
		api.ReportPing(backendUrl, ip, 0, containerStatusDOWN)
	}

	pinger, err := probing.NewPinger(ip)

	if err != nil {
		reportDown()
		return
	}

	pinger.Count = 1
	pinger.Timeout = time.Duration(pingPeriodSec) * time.Second
	err = pinger.Run()
	stats := pinger.Statistics()

	if err != nil || stats.PacketsRecv == 0 {
		reportDown()
		return
	}

	pingTimeMs := int32(stats.AvgRtt.Milliseconds())
	api.ReportPing(backendUrl, ip, pingTimeMs, containerStatusUP)
}
