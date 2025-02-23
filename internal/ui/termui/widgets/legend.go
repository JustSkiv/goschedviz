// Package widgets provides terminal UI components using termui library.
package widgets

import (
	tui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// PlotLegend displays color information for the plot lines
type PlotLegend struct {
	*widgets.Paragraph
}

// NewPlotLegend creates a new legend widget
func NewPlotLegend() *PlotLegend {
	l := &PlotLegend{
		Paragraph: widgets.NewParagraph(),
	}
	l.Title = "Legend"
	l.BorderStyle.Fg = tui.ColorWhite

	// Компактное отображение с цветными линиями и текстом
	l.Text = "[-- [GRQ]](fg:green)\n" +
		"[-- [LRQ]](fg:magenta)\n" +
		"[-- [THR]](fg:red)\n" +
		"[-- [IDL]](fg:yellow)\n" +
		"[-- [GRT]](fg:cyan)"

	return l
}
