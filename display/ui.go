package display

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/gizak/termui/v3"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/hugo-sv/webmonitor/statistics"
)

// View is a Structure containing all relevant informatiosn to display the UI
type View struct {
	// Monitored URLS
	Urls []string
	// Associated Statistics
	URLStatistics map[string][3]*statistics.Statistic
	// The user's active timeframe
	ActiveTimeframe int
	// String representation for ActiveTimeframe
	TimeframeRepr map[int]string
	// The user's active Detailed view
	ActiveWebsite int
	// Alerts Messages
	AlertMessages []string
	// Alert Scroll position
	AlertOffset int
}

// Init initialize the UI
func Init() <-chan termui.Event {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	return ui.PollEvents()
}

// Close Closes the UI
func Close() {
	ui.Close()
}

// RenderLayout renders the UI Layout
func RenderLayout(uiView View) {
	p := widgets.NewParagraph()
	p.Title = " Webmonitor "
	p.Text = "Website monitoring tool. Press q to quit."
	p.TextStyle.Fg = ui.ColorYellow
	p.SetRect(0, 0, 75, 3)
	p.BorderStyle.Fg = ui.ColorCyan
	ui.Render(p)
	// Render Sub-layouts
	renderAlertsLayout()
	renderStatisticsLayout(uiView)
}

// renderAlertsLayout renders the Alerts Layout
func renderAlertsLayout() {
	p := widgets.NewParagraph()
	p.Title = " Alerts "
	p.Text = "Press up and down to scroll through alerts.\n\n There are no alerts."
	p.TextStyle.Fg = ui.ColorYellow
	p.SetRect(75, 0, 150, 50)
	p.BorderStyle.Fg = ui.ColorCyan
	ui.Render(p)
}

// renderStatisticsLayout renders the Statistics Layout
func renderStatisticsLayout(uiView View) {
	p1 := widgets.NewParagraph()
	p1.Title = fmt.Sprintf(" Statistics : last %v ", uiView.TimeframeRepr[uiView.ActiveTimeframe])
	p1.Text = fmt.Sprintf("Press s to switch to a %v timeframe.\n\n Statistics are loading ...", uiView.TimeframeRepr[3-uiView.ActiveTimeframe])
	p1.TextStyle.Fg = ui.ColorYellow
	p1.SetRect(0, 3, 75, 26)
	p1.BorderStyle.Fg = ui.ColorCyan

	p2 := widgets.NewParagraph()
	p2.Title = fmt.Sprintf(" Details %v ", uiView.ActiveWebsite)
	p2.Text = "Press a website's id to view details"
	p2.TextStyle.Fg = ui.ColorYellow
	p2.SetRect(0, 26, 75, 50)
	p2.BorderStyle.Fg = ui.ColorCyan
	ui.Render(p1, p2)
}

// renderStatTable renders a table of statistics
func renderStatTable(Table [][]string) {
	g := widgets.NewTable()
	g.SetRect(0, 5, 75, 26)
	g.TextAlignment = ui.AlignCenter
	g.ColumnWidths = []int{6, 34, 10, 10, 10}
	g.Rows = Table
	ui.Render(g)
}

// RenderAlerts renders a list of alert messages
func RenderAlerts(uiView View) {
	p := widgets.NewParagraph()
	p.Text = strings.Join(uiView.AlertMessages[uiView.AlertOffset:], "\n")
	p.SetRect(75, 2, 150, 50)
	p.TextStyle.Fg = ui.ColorWhite
	ui.Render(p)
}

// RenderStats isolate the relevant Satistics to display them.
func RenderStats(uiView View) {
	renderStatisticsLayout(uiView)

	// Processing the general stats table headers
	statsHeaders := []string{
		"Id",
		"Website",
		"Avg",
		"Max",
		"Availability",
	}
	Table := [][]string{statsHeaders}
	// For each URL
	var urlStatistic *statistics.Statistic
	for id, url := range uiView.Urls {
		urlStatistic = uiView.URLStatistics[url][uiView.ActiveTimeframe]
		if !math.IsNaN(urlStatistic.Average()) {
			// Append Statistics
			Table = append(Table, []string{
				fmt.Sprint(id),
				Shorten(url),
				fmt.Sprintf("%.0f", urlStatistic.Average()),
				fmt.Sprintf("%v", urlStatistic.MaxResponseTime()),
				fmt.Sprintf("%.0f%%", urlStatistic.Availability()*100.0),
			})
		}
	}
	// Rendering the updated table
	renderStatTable(Table)

	// Processing the detailed view
	detailedStatistics := uiView.URLStatistics[uiView.Urls[uiView.ActiveWebsite]]
	detailedHeaders := []string{
		"TimeFrame",
		"Avg",
		"Max",
		"Availability",
		"Codes",
	}
	detailTable := [][]string{detailedHeaders}
	for id, statistic := range detailedStatistics {
		if !math.IsNaN(urlStatistic.Average()) {
			// Append Statistics
			detailTable = append(detailTable, []string{
				uiView.TimeframeRepr[id],
				fmt.Sprintf("%.0f", statistic.Average()),
				fmt.Sprintf("%v", statistic.MaxResponseTime()),
				fmt.Sprintf("%.0f%%", statistic.Availability()*100.0),
				StatusCodeMapToString(statistic.StatusCodeCount),
			})
		}
	}
	// Processing the sparkline graph
	plotValues := detailedStatistics[2].RecentResponseTime()

	// Rendering the detailed view
	renderStatDetails(uiView, detailTable, plotValues)
}

// renderStatDetails renders detailed statistics view
func renderStatDetails(uiView View, detailTable [][]string, plotValues []float64) {
	// Detailed table
	g := widgets.NewTable()
	g.SetRect(0, 28, 75, 37)
	g.TextAlignment = ui.AlignCenter
	g.ColumnWidths = []int{10, 10, 10, 10, 30}
	g.Rows = detailTable

	// Detailed sparkline
	slc := widgets.NewSparkline()
	slc.Data = plotValues
	slc.LineColor = ui.ColorGreen

	lc := widgets.NewSparklineGroup(slc)
	lc.Title = " Recent Response Time "

	lc.SetRect(1, 37, 75, 49)

	ui.Render(g, lc)
}
