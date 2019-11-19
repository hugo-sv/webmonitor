package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

// JSONInput struct which contains an array of websites
type JSONInput struct {
	Timeout  int       `json:"timeout"`
	Websites []Website `json:"websites"`
}

// Website struct which contains an url and an interval
type Website struct {
	URL      string `json:"url"`
	Interval int    `json:"interval"`
}

// ParseFlags parse and returns the flags of the webmonitor cli command : Websites, Check interval, and timeout
func ParseFlags() (int, []string, map[string]int) {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("No JSON input file path specified")
		return 0, make([]string, 0), make(map[string]int)
	}
	jsonpath := flag.Args()[0]
	// Open the file
	jsonFile, err := os.Open(jsonpath)
	if err != nil {
		// Handle error
		fmt.Println(err)
		return 0, make([]string, 0), make(map[string]int)
	}
	defer jsonFile.Close()
	// Read the Json
	var input JSONInput
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &input)

	// Parsing the JSON
	intervals := make(map[string]int)
	urls := make([]string, 0)
	for _, website := range input.Websites {
		// Interval should be greater than 1, no duplicate URL
		if website.Interval >= 1 && intervals[website.URL] == 0 {
			intervals[website.URL] = website.Interval
			urls = append(urls, website.URL)
		}
	}
	// If timeout invalid
	if input.Timeout <= 1 {
		fmt.Println("Timeout specified in JSON should be an integer greater than 1")
		return 0, make([]string, 0), make(map[string]int)
	}

	return input.Timeout, urls, intervals
}
