package rollbar

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

// A Widget represents a Rollbar widget
type Widget struct {
	view.KeyboardWidget
	view.ScrollableWidget
	view.StatableWidget

	items    *Result
	settings *Settings
}

// NewWidget creates a new instance of a widget
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
	if widget.Disabled() {
		return
	}

	items, err := CurrentActiveItems(
		widget.settings.accessToken,
		widget.settings.assignedToName,
		widget.settings.activeOnly,
	)

	if err != nil {
		widget.Redraw(widget.CommonSettings().Title, err.Error(), true)
		return
	}
	widget.items = &items.Results
	widget.SetItemCount(len(widget.items.Items))

	widget.Render()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) Render() {
	if widget.items == nil {
		return
	}

	title := fmt.Sprintf("%s - %s", widget.CommonSettings().Title, widget.settings.projectName)
	widget.Redraw(title, widget.contentFrom(widget.items), false)
}

func (widget *Widget) contentFrom(result *Result) string {
	if result == nil {
		return "No results"
	}
	var str string
	if len(result.Items) > widget.settings.count {
		result.Items = result.Items[:widget.settings.count]
	}
	for idx, item := range result.Items {

		row := fmt.Sprintf(
			"[%s] [%s] %s [%s] %s [%s]count: %d [%s]%s",
			widget.RowColor(idx),
			levelColor(&item),
			item.Level,
			statusColor(&item),
			item.Title,
			widget.RowColor(idx),
			item.TotalOccurrences,
			"blue",
			item.Environment,
		)
		str += utils.HighlightableHelper(widget.View, row, idx, len(item.Title))
	}

	return str
}

func statusColor(item *Item) string {
	switch item.Status {
	case "active":
		return "red"
	case "resolved":
		return "green"
	default:
		return "red"
	}
}
func levelColor(item *Item) string {
	switch item.Level {
	case "error":
		return "red"
	case "critical":
		return "green"
	case "warning":
		return "yellow"
	default:
		return "grey"
	}
}

func (widget *Widget) openBuild() {
	if widget.GetSelected() >= 0 && widget.items != nil && widget.GetSelected() < len(widget.items.Items) {
		item := &widget.items.Items[widget.GetSelected()]

		utils.OpenFile(
			fmt.Sprintf(
				"https://rollbar.com/%s/%s/%s/%d",
				widget.settings.projectOwner,
				widget.settings.projectName,
				"items",
				item.ID,
			),
		)
	}
}
