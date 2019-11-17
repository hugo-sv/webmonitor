package main

import (
	"fmt"
	"math"
	"time"

	"github.com/hugo-sv/webmonitor/cli"
	"github.com/hugo-sv/webmonitor/display"
	"github.com/hugo-sv/webmonitor/monitor"
	"github.com/hugo-sv/webmonitor/statistics"
)

// CheckStats is a data structure to store results form the monitoring CheckWithTimeOut function
type CheckStats struct {
	url          string
	responseTime int
	StatusCode   int
}

// DisplayStats isolate the relevant Satistics to display them.
func DisplayStats(urls []string, urlStatistics map[string][3]*statistics.Statistic, activeView int) {
	// Setting up stats table headers
	statsHeaders := []string{
		"Website",
		"Avg",
		"Max",
		"Availability",
		"Codes",
	}
	Table := [][]string{statsHeaders}
	// For each URL
	var urlStatistic *statistics.Statistic
	for _, url := range urls {
		urlStatistic = urlStatistics[url][activeView]
		// Append Statistics
		Table = append(Table, []string{
			display.Shorten(url),
			fmt.Sprintf("%.0f", urlStatistic.Average()),
			fmt.Sprintf("%v", urlStatistic.MaxResponseTime()),
			fmt.Sprintf("%.0f%%", urlStatistic.Availability()*100.0),
			display.StatusCodeMapToString(urlStatistic.StatusCodeCount()),
		})
	}
	// Render the update table
	display.RenderStatTable(Table)
}

func main() {
	// Retrieving the command's flags
	checkInterval, timeout, urls := cli.ParseFlags()
	// Setting up the Statistics and alerts objects
	AlertMessages := []string{}
	AlertStartMessage := "[%v] Website %s is down. Availability : %.0f %%"
	AlertEndMessage := "[%v] Website %s is up."

	urlStatistics := make(map[string][3]*statistics.Statistic)
	for _, url := range urls {
		// Keeping track of Statistics for 2min, 10min and 1h timeframes
		urlStatistics[url] = [3]*statistics.Statistic{
			statistics.NewStatistic(int(math.Ceil(float64(2*60) / float64(checkInterval)))),
			statistics.NewStatistic(int(math.Ceil(float64(10*60) / float64(checkInterval)))),
			statistics.NewStatistic(int(math.Ceil(float64(60*60) / float64(checkInterval)))),
		}
	}
	// Fetch Ticker will trigger the Check
	fetchTicker := time.NewTicker(time.Second * time.Duration(checkInterval))
	defer fetchTicker.Stop()
	// These Display tickers will trigger the display of stats for the past 10min and 1h respectively
	displayTicker1 := time.NewTicker(time.Second * time.Duration(10))
	defer displayTicker1.Stop()
	displayTicker2 := time.NewTicker(time.Second * time.Duration(60))
	defer displayTicker2.Stop()
	// The user's active view
	// 1 means the users have the past 10min stats displayed
	// 2 means the users have the past 1h stats displayed
	activeView := 1
	// Setting up the UI display
	uiEvents := display.Init()
	defer display.Close()
	display.RenderLayout(activeView)
	// Listening to tickers and UI Events ...
	var previousAvailability, currentAvailability float64
	var responseTime, StatusCode int
	statsToDisplay := false
	statsMessage := make(chan CheckStats)
	for {
		select {
		// Check Ticker
		case <-fetchTicker.C:
			for _, url := range urls {
				// Checking the website within a goroutine
				go func(url string) {
					responseTime, StatusCode := monitor.CheckWithTimeout(url, timeout)
					statsMessage <- CheckStats{url, responseTime, StatusCode}
				}(url)
			}
		// Catching the result of a Check operation
		case Stats := <-statsMessage:
			responseTime = Stats.responseTime
			StatusCode = Stats.StatusCode
			// Pulling the previous availability
			previousAvailability = urlStatistics[Stats.url][0].Availability()
			for i, urlStatistic := range urlStatistics[Stats.url] {
				// Updating the records
				urlStatistic.AddRecord(responseTime, StatusCode)
				statsToDisplay = true
				// Handeling alerts with the 2 min timeframe stats
				if i == 0 {
					currentAvailability = urlStatistic.Availability()
					// If 80% threshold is crossed, or website is unavailable from the start
					if currentAvailability < 0.8 && (previousAvailability >= 0.8 || math.IsNaN(previousAvailability)) {
						// Append the alert
						AlertMessages = append(AlertMessages,
							fmt.Sprintf(AlertStartMessage,
								time.Now().Format(time.Kitchen),
								display.Shorten(Stats.url),
								currentAvailability*100.0,
							))
						// Update the UI
						display.RenderAlerts(AlertMessages)
					}
					// If availability is back above the 80% threshold
					if currentAvailability > 0.8 && previousAvailability <= 0.8 {
						// Append the alert
						AlertMessages = append(AlertMessages,
							fmt.Sprintf(AlertEndMessage,
								time.Now().Format(time.Kitchen),
								display.Shorten(Stats.url),
							))
						// Update the UI
						display.RenderAlerts(AlertMessages)
					}
				}
			}
		// 10 min display Ticker
		case <-displayTicker1.C:
			if activeView == 1 {
				DisplayStats(urls, urlStatistics, activeView)
			}
		// 1 h display Ticker
		case <-displayTicker2.C:
			if activeView == 2 {
				DisplayStats(urls, urlStatistics, activeView)
			}
		// UI events
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "c":
				// Clearing the alerts
				AlertMessages = []string{}
				display.RenderAlertsLayout()
			case "s":
				// Switching the active view
				activeView = 3 - activeView
				display.RenderStatisticsLayout(activeView)
				if statsToDisplay {
					// Displaying stats table if and only if there are data to display
					DisplayStats(urls, urlStatistics, activeView)
				}
			}
		}
	}
}
