package monitor

import (
	"net/http"
	"time"
)

// CheckWithTimeout Checks a website, and returns the current stats of a website (response time and response code), unless it times out.
func CheckWithTimeout(url string, timeout int) (int, int) {
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
		return responseTime, 408
	}
	defer resp.Body.Close()
	return responseTime, resp.StatusCode
}
