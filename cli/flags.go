package cli

import "flag"

// ParseFlags returns the flags of the webmonitor command : Check interval, timeout and URLs
func ParseFlags() (int, int, []string) {
	var checkInterval int
	flag.IntVar(&checkInterval, "interval", 5, "In seconds, interval at which to check the websites availability")
	var timeout int
	flag.IntVar(&timeout, "timeout", 10, "In seconds, timeout at which to stop checking the websites availability")
	flag.Parse()
	urls := flag.Args()
	return checkInterval, timeout, urls
}
