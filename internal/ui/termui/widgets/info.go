package widgets

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui/v3"
	twidgets "github.com/gizak/termui/v3/widgets"
)

// InfoPanel shows monitoring statistics and help
type InfoPanel struct {
	*twidgets.Paragraph
}

// NewInfoPanel creates a new info panel widget
func NewInfoPanel() *InfoPanel {
	p := &InfoPanel{
		Paragraph: twidgets.NewParagraph(),
	}
	p.Title = "Information"
	p.BorderStyle.Fg = ui.ColorCyan
	return p
}

// Update refreshes info panel with current statistics
func (p *InfoPanel) Update(historyLen, maxGRQ, maxLRQ int) {
	p.Text = fmt.Sprintf(
		"Press 'q' to exit\n"+
			"Last update: %s\n"+
			"History points: %d\n"+
			"Max GRQ: %d\n"+
			"Max LRQ (sum): %d\n",
		time.Now().Format("15:04:05"),
		historyLen,
		maxGRQ,
		maxLRQ,
	)
}
