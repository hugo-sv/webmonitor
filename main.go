package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/hugo-sv/webmonitor/cli"
	"github.com/hugo-sv/webmonitor/display"
	"github.com/hugo-sv/webmonitor/monitor"
	"github.com/hugo-sv/webmonitor/statistics"
)

func main() {
	// Retrieving the cli command's flags
	timeout, urls, intervals, uiEnabled := cli.ParseFlags()
	if len(urls) == 0 {
		// There are no URL to track
		return
	}
	// Channel messages
	stop := make(chan struct{}, len(urls))
	statsMessage := make(chan monitor.CheckStats)
	// Setting up the Statistics system
	urlStatistics := make(map[string][3]*statistics.Statistic)
	var checkInterval int
	for _, url := range urls {
		checkInterval = intervals[url]
		// Keeping track of enough records for 2min, 10min and 1h timeframes
		urlStatistics[url] = [3]*statistics.Statistic{
			statistics.NewStatistic(int(math.Ceil(float64(2*60) / float64(checkInterval)))),
			statistics.NewStatistic(int(math.Ceil(float64(10*60) / float64(checkInterval)))),
			statistics.NewStatistic(int(math.Ceil(float64(60*60) / float64(checkInterval)))),
		}
		// Starting a goroutine fetching data for this URL
		go monitor.CheckOnTicks(url, checkInterval, timeout, stop, statsMessage)
	}
	// These Display tickers will refresh the stats display every 10sec and 1min for the past 10min and 1h respectively
	displayTicker1 := time.NewTicker(time.Second * time.Duration(10))
	defer displayTicker1.Stop()
	displayTicker2 := time.NewTicker(time.Second * time.Duration(60))
	defer displayTicker2.Stop()
	// Setting up the UI display
	uiView := display.View{
		UIEnabled:       uiEnabled,
		Urls:            urls,
		URLStatistics:   urlStatistics,
		TimeframeRepr:   map[int]string{0: "2min", 1: "10min", 2: "1h"},
		ActiveWebsite:   0,
		ActiveTimeframe: 1,
		AlertMessages:   []string{},
		AlertOffset:     0,
	}
	uiEvents := display.Init(uiView)
	defer display.Close(uiView)
	display.RenderLayout(uiView)
	// Listening to tickers and UI Events ...
	var previousAvailability, currentAvailability float64
	for {
		select {
		// Catching the result of a Check operation
		case stats := <-statsMessage:
			// Pulling the previous availability
			previousAvailability = urlStatistics[stats.URL][0].Availability()
			for i, urlStatistic := range urlStatistics[stats.URL] {
				// Updating the records
				urlStatistic.AddRecord(stats.ResponseTime, stats.StatusCode)
				// Handeling alerts with the 2 min timeframe stats
				if i == 0 {
					currentAvailability = urlStatistic.Availability()
					// If 80% threshold is crossed, or website is unavailable from the start
					if currentAvailability < 0.8 && (previousAvailability >= 0.8 || math.IsNaN(previousAvailability)) {
						// Append the alert
						uiView.AlertMessages = append(uiView.AlertMessages,
							fmt.Sprintf("Website %s is down. availability=%.0f %%, time=%v",
								display.Shorten(stats.URL),
								currentAvailability*100.0,
								time.Now().Format(time.Kitchen),
							))
						// Update the UI
						go display.RenderAlerts(uiView)
					}
					// If availability is back above the 80% threshold
					if currentAvailability >= 0.8 && previousAvailability < 0.8 {
						// Append the alert
						uiView.AlertMessages = append(uiView.AlertMessages,
							fmt.Sprintf("Website %s is up, time=%v",
								display.Shorten(stats.URL),
								time.Now().Format(time.Kitchen),
							))
						// Update the UI
						go display.RenderAlerts(uiView)
					}
				}
			}
		// 10 min display Ticker
		case <-displayTicker1.C:
			go display.RenderStats(uiView, 1)
		// 1 h display Ticker
		case <-displayTicker2.C:
			go display.RenderStats(uiView, 2)
		// UI events
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				// Sending as many stop messages as there are Fetch goroutines running
				for range urls {
					stop <- struct{}{}
				}
				return
			case "<Up>":
				// Scrolling up
				if uiView.AlertOffset > 0 {
					uiView.AlertOffset--
				}
				go display.RenderAlerts(uiView)
			case "<Down>":
				// Scrolling down
				if uiView.AlertOffset < len(uiView.AlertMessages)-1 {
					uiView.AlertOffset++
				}
				go display.RenderAlerts(uiView)
			case "s":
				// Switching the active view between 1 and 2
				uiView.ActiveTimeframe = 3 - uiView.ActiveTimeframe
				go display.RenderStats(uiView, uiView.ActiveTimeframe)
			}
			// If the pressed key is a number within the nuber of websites' range
			v, err := strconv.Atoi(e.ID)
			if err == nil && v < len(urls) {
				uiView.ActiveWebsite = v
				// Updating Statistics layout and views
				go display.RenderStats(uiView, uiView.ActiveTimeframe)
			}
		}
	}
}
