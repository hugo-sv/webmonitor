package display

import (
	"log"
	"strings"

	"github.com/gizak/termui/v3"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

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
func RenderLayout(ActiveTimeframe int) {
	p := widgets.NewParagraph()
	p.Title = " Webmonitor "
	p.Text = "Website monitoring tool. Press q to quit."
	p.SetRect(0, 0, 75, 3)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan
	ui.Render(p)
	// Render sub layouts
	RenderStatisticsLayout(ActiveTimeframe)
	RenderAlertsLayout()
}

// RenderStatisticsLayout renders the Statistics Layout
func RenderStatisticsLayout(ActiveTimeframe int) {
	p := widgets.NewParagraph()
	p.Title = " Statistics : last "
	p.Text = "Press s to switch to a "
	switch ActiveTimeframe {
	case 1:
		p.Title += "10 min "
		p.Text += "1 h"
	default:
		p.Title += "1 h "
		p.Text += "10 min"
	}
	p.Text += " timeframe.\n\n Statistics are loading ..."
	p.SetRect(75, 0, 150, 50)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan
	ui.Render(p)
}

// RenderAlertsLayout renders the Alerts Layout
func RenderAlertsLayout() {
	p := widgets.NewParagraph()
	p.Title = " Alerts "
	p.Text = "Press c to clear alerts.\n\n There are no alerts."
	p.SetRect(0, 3, 75, 50)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan
	ui.Render(p)
}

// RenderStatTable renders a table of statistics
func RenderStatTable(Table [][]string) {
	g := widgets.NewTable()
	g.SetRect(75, 2, 150, 50)
	g.TextAlignment = ui.AlignCenter
	g.ColumnWidths = []int{30, 6, 6, 6, 22}
	g.Rows = Table
	ui.Render(g)
}

// RenderAlerts renders a list of alert messages
func RenderAlerts(Messages []string) {
	p := widgets.NewParagraph()
	p.Text = strings.Join(Messages, "\n")
	p.SetRect(0, 5, 75, 50)
	p.TextStyle.Fg = ui.ColorWhite
	ui.Render(p)
}
