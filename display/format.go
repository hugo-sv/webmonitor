package display

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Shorten shortens a URL by removing non relevant headers.
func Shorten(url string) string {
	r, _ := regexp.Compile("^(http(s)?://)?(www.)?")
	return r.ReplaceAllString(url, "")
}

// StatusCodeMapToString convert the Status Code map to a displayable String of Codes and count, sorted by count.
func StatusCodeMapToString(StatusCode map[int]int) string {
	// statusCodeCount is structure used to sort a StatusCode count Map.
	type statusCodeCount struct {
		StatusCode int
		Count      int
	}
	// Retrieving only Status counted at least once
	Statuses := make([]statusCodeCount, 0)
	for statusCode, count := range StatusCode {
		if count > 0 {
			Statuses = append(Statuses, statusCodeCount{statusCode, count})
		}
	}
	// Sorting ths Statuses
	sort.Slice(Statuses, func(i, j int) bool {
		return Statuses[i].Count > Statuses[j].Count
	})
	// Building string representation
	repr := ""
	for _, statusCodeCount := range Statuses {
		repr += fmt.Sprintf("%v:%v, ", statusCodeCount.StatusCode, statusCodeCount.Count)
	}
	return strings.TrimSuffix(repr, ", ")
}
