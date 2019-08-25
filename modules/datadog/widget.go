package datadog

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
	datadog "github.com/zorkian/go-datadog-api"
)

type Widget struct {
	view.KeyboardWidget
	view.ScrollableWidget
	view.StatableWidget

	monitors []datadog.Monitor
	settings *Settings
}

func NewWidget(app *tview.Application, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		KeyboardWidget:   view.NewKeyboardWidget(app, pages, settings.common),
		ScrollableWidget: view.NewScrollableWidget(app, settings.common, true),

		settings: settings,
	}

	widget.SetRenderFunction(widget.Render)
	widget.initializeKeyboardControls()
	widget.View.SetInputCapture(widget.InputCapture)

	widget.KeyboardWidget.SetView(widget.View)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	monitors, monitorErr := widget.Monitors()

	if monitorErr != nil {
		widget.monitors = nil
		widget.SetItemCount(0)
		widget.Redraw(widget.CommonSettings().Title, monitorErr.Error(), true)
		return
	}
	triggeredMonitors := []datadog.Monitor{}

	for _, monitor := range monitors {
		state := *monitor.OverallState
		if state == "Alert" {
			triggeredMonitors = append(triggeredMonitors, monitor)
		}
	}
	widget.monitors = triggeredMonitors
	widget.SetItemCount(len(widget.monitors))

	widget.Render()
}

func (widget *Widget) Render() {
	content := widget.contentFrom(widget.monitors)
	widget.Redraw(widget.CommonSettings().Title, content, false)
}

func (widget *Widget) HelpText() string {
	return widget.KeyboardWidget.HelpText()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) contentFrom(triggeredMonitors []datadog.Monitor) string {
	var str string

	if len(triggeredMonitors) > 0 {
		str += fmt.Sprintf(
			" %s\n",
			"[red]Triggered Monitors[white]",
		)
		for idx, triggeredMonitor := range triggeredMonitors {
			row := fmt.Sprintf(`[%s][red] %s[%s]`,
				widget.RowColor(idx),
				*triggeredMonitor.Name,
				widget.RowColor(idx),
			)
			str += utils.HighlightableHelper(widget.View, row, idx, len(*triggeredMonitor.Name))
		}
	} else {
		str += fmt.Sprintf(
			" %s\n",
			"[green]No Triggered Monitors[white]",
		)
	}

	return str
}

func (widget *Widget) openItem() {

	sel := widget.GetSelected()
	if sel >= 0 && widget.monitors != nil && sel < len(widget.monitors) {
		item := &widget.monitors[sel]
		utils.OpenFile(fmt.Sprintf("https://app.datadoghq.com/monitors/%d?q=*", *item.Id))
	}
}
