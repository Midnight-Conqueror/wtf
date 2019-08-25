package gitter

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

// A Widget represents a Gitter widget
type Widget struct {
	view.KeyboardWidget
	view.ScrollableWidget
	view.StatableWidget

	messages []Message
	settings *Settings
}

// NewWidget creates a new instance of a widget
func NewWidget(app *tview.Application, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		KeyboardWidget:   view.NewKeyboardWidget(app, pages, settings.common),
		ScrollableWidget: view.NewScrollableWidget(app, settings.common, true),

		settings: settings,
	}

	widget.SetRenderFunction(widget.Refresh)
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

	room, err := GetRoom(widget.settings.roomURI, widget.settings.apiToken)
	if err != nil {
		widget.Redraw(widget.CommonSettings().Title, err.Error(), true)
		return
	}

	if room == nil {
		widget.Redraw(widget.CommonSettings().Title, "No room", true)
		return
	}

	messages, err := GetMessages(room.ID, widget.settings.numberOfMessages, widget.settings.apiToken)

	if err != nil {
		widget.Redraw(widget.CommonSettings().Title, err.Error(), true)
		return
	}
	widget.messages = messages
	widget.SetItemCount(len(messages))

	widget.display()
}

func (widget *Widget) HelpText() string {
	return widget.KeyboardWidget.HelpText()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) display() {
	if widget.messages == nil {
		return
	}

	title := fmt.Sprintf("%s - %s", widget.CommonSettings().Title, widget.settings.roomURI)

	widget.Redraw(title, widget.contentFrom(widget.messages), true)
}

func (widget *Widget) contentFrom(messages []Message) string {
	var str string
	for idx, message := range messages {
		row := fmt.Sprintf(
			`[%s] [blue]%s [lightslategray]%s: [%s]%s [aqua]%s`,
			widget.RowColor(idx),
			message.From.DisplayName,
			message.From.Username,
			widget.RowColor(idx),
			message.Text,
			message.Sent.Format("Jan 02, 15:04 MST"),
		)

		str += utils.HighlightableHelper(widget.View, row, idx, len(message.Text))
	}

	return str
}

func (widget *Widget) openMessage() {
	sel := widget.GetSelected()
	if sel >= 0 && widget.messages != nil && sel < len(widget.messages) {
		message := &widget.messages[sel]
		utils.OpenFile(message.Text)
	}
}
