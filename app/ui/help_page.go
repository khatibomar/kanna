package ui

import (
	"fmt"

	"codeberg.org/omarkhatib/kanna/app/core"
	"codeberg.org/omarkhatib/kanna/app/ui/utils"
	"github.com/rivo/tview"
)

const (
	padding = 20
)

// HelpPage : This struct contains the grid for the help page.
type HelpPage struct {
	Grid *tview.Grid
	Core *core.Kanna
}

// ShowHelpPage : Make the app show the help page.
func ShowHelpPage(core *core.Kanna) {
	helpPage := newHelpPage(core)

	core.TView.SetFocus(helpPage.Grid)
	core.PageHolder.AddPage(utils.HelpPageID, helpPage.Grid, true, true)
}

// newHelpPage : Creates a new help page.
func newHelpPage(core *core.Kanna) *HelpPage {
	formatString := fmt.Sprintf("%%-%ds:%%%ds\n", padding, padding)
	// Set up the help text.
	helpText := "Keyboard Mappings\n" +
		"-----------------------------\n\n" +
		"Universal\n" +
		fmt.Sprintf(formatString, "Ctrl + K", "Keybinds/Help") +
		fmt.Sprintf(formatString, "Ctrl + S", "Search") +
		"\nAnime Page\n" +
		fmt.Sprintf(formatString, "Ctrl + E", "Select mult.") +
		fmt.Sprintf(formatString, "Ctrl + A", "Toggle All") +
		fmt.Sprintf(formatString, "Ctrl + R", "Toggle Watched Status") +
		fmt.Sprintf(formatString, "Enter", "Start download") +
		"\nOthers\n" +
		fmt.Sprintf(formatString, "Esc", "Go back") +
		fmt.Sprintf(formatString, "Ctrl + F/B", "Next/Prev Page") +
		"\n"

	// Create TextView to show the help information.
	help := tview.NewTextView()
	// Set TextView attributes.
	help.SetText(helpText).
		SetTextAlign(tview.AlignCenter).
		SetBorderColor(utils.HelpPageBorderColor).
		SetBorder(true)

	// Create a new grid for the text view, so we can align it to the center.
	dimensions := []int{-1, -1, -1, -1, -1, -1}
	grid := utils.NewGrid(dimensions, dimensions).
		AddItem(help, 0, 0, 6, 6, 0, 0, true)

	helpPage := &HelpPage{
		Grid: grid,
		Core: core,
	}
	helpPage.setHandlers()

	return helpPage
}
