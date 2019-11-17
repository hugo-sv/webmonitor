package display

import (
	"encoding/json"
	"log"
	"regexp"
)

// Shorten shortens a URL by removing url headers.
func Shorten(url string) string {
	r, _ := regexp.Compile("^(http(s)?://)?(www.)?")
	return r.ReplaceAllString(url, "")
}

// StatusCodeMapToString convert the Status Code map to a displayable String
func StatusCodeMapToString(StatusCode map[int]int) string {
	empData, err := json.Marshal(StatusCode)
	if err != nil {
		log.Fatalf("failed to convert map to string : %v", err)
		return ""
	}
	return string(empData)
}
