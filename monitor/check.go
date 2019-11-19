package monitor

import (
	"net/http"
	"time"
)

// CheckStats is a data structure to store and send results from the CheckWithTimeout function
type CheckStats struct {
	URL          string
	ResponseTime int
	StatusCode   int
}

// CheckWithTimeout Checks a website, and returns the current response time and response code of a website, unless it times out.
func CheckWithTimeout(url string, timeout int) CheckStats {
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	t := time.Now()
	// Checking website
	resp, err := client.Get(url)
	// Computing Total response time
	responseTime := int(time.Now().Sub(t).Milliseconds())
	// If there are no response, or a timeout
	if err != nil {
		return CheckStats{url, responseTime, 408}
	}
	defer resp.Body.Close()
	return CheckStats{url, responseTime, resp.StatusCode}
}

// CheckOnTicks Regularly checks a website, and send back the stats as a channel message.
func CheckOnTicks(url string, checkInterval int, timeout int, stop chan struct{}, statsMessage chan CheckStats) {
	// Data will be fetched at every checkInterval
	fetchTicker := time.NewTicker(time.Second * time.Duration(checkInterval))
	defer fetchTicker.Stop()
	for {
		select {
		// Message to stop the goroutine and ticker
		case <-stop:
			return
		// Fetch Ticker
		case <-fetchTicker.C:
			// Checking the website within a goroutine and send back the results
			go func() {
				statsMessage <- CheckWithTimeout(url, timeout)
			}()
		}
	}
}
