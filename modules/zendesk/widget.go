package zendesk

import (
	"fmt"
	"log"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

// A Widget represents a Zendesk widget
type Widget struct {
	view.KeyboardWidget
	view.ScrollableWidget
	view.StatableWidget

	result   *TicketArray
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
	ticketArray, err := widget.newTickets()
	ticketArray.Count = len(ticketArray.Tickets)
	if err != nil {
		log.Fatal(err)
	} else {
		widget.result = ticketArray
	}

	widget.Render()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) Render() {
	title := fmt.Sprintf("%s (%d)", widget.CommonSettings().Title, widget.result.Count)
	widget.Redraw(title, widget.textContent(widget.result.Tickets), false)
}

func (widget *Widget) textContent(items []Ticket) string {
	if len(items) == 0 {
		return fmt.Sprintf("No unassigned tickets in queue - woop!!")
	}

	str := ""
	for idx, data := range items {
		str += widget.format(data, idx)
	}

	return str
}

func (widget *Widget) format(ticket Ticket, idx int) string {
	textColor := widget.settings.common.Colors.Background
	if idx == widget.GetSelected() {
		textColor = widget.settings.common.Colors.BorderFocused
	}

	requesterName := widget.parseRequester(ticket)
	str := fmt.Sprintf(" [%s:]%d - %s\n %s\n\n", textColor, ticket.Id, requesterName, ticket.Subject)
	return str
}

// this is a nasty means of extracting the actual name of the requester from the Via interface of the Ticket.
// very very open to improvements on this
func (widget *Widget) parseRequester(ticket Ticket) interface{} {
	viaMap := ticket.Via
	via := viaMap.(map[string]interface{})
	source := via["source"]
	fromMap, _ := source.(map[string]interface{})
	from := fromMap["from"]
	fromValMap := from.(map[string]interface{})
	fromName := fromValMap["name"]
	return fromName
}

func (widget *Widget) openTicket() {
	sel := widget.GetSelected()
	if sel >= 0 && widget.result != nil && sel < len(widget.result.Tickets) {
		issue := &widget.result.Tickets[sel]
		ticketURL := fmt.Sprintf("https://%s.zendesk.com/agent/tickets/%d", widget.settings.subdomain, issue.Id)
		utils.OpenFile(ticketURL)
	}
}
