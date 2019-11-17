package monitor

import "testing"

func TestCheckStatusCode(t *testing.T) {
	// Test if the proper status codes are returned
	cases := []struct {
		in   string
		want int
	}{
		{"https://httpstat.us/304", 304},
		{"https://httpstat.us/512", 512},
		{"https://httpstat.us/200", 200},
		{"https://httpstat.us/200?sleep=5100", 408},
		{"https://jehbqkqbwelkjsd.com/", 408},
	}
	for _, c := range cases {
		_, got := CheckWithTimeout(c.in, 5)
		if got != c.want {
			t.Errorf("Reverse(%q) == %v, want %v", c.in, got, c.want)
		}
	}
}

func TestCheckResponseTime(t *testing.T) {
	// With httpstat.us, Test if ResponseTime is not lower than expected
	cases := []struct {
		in   string
		want int
	}{
		{"https://httpstat.us/200?sleep=10", 10},
		{"https://httpstat.us/200?sleep=100", 100},
		{"https://httpstat.us/200?sleep=1000", 1000},
		{"https://httpstat.us/200?sleep=10000", 5000},
	}
	for _, c := range cases {
		got, _ := CheckWithTimeout(c.in, 5)
		if got < c.want {
			t.Errorf("Reverse(%q) < %v, want %v", c.in, got, c.want)
		}
	}
}
