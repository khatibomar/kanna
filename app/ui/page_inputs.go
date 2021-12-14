package ui

import (
	"context"
	"log"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/khatibomar/angoslayer"
	"github.com/khatibomar/tkanna/app/core"
	"github.com/khatibomar/tkanna/app/ui/utils"
)

// SetUniversalHandlers : Set universal inputs for the app.
func SetUniversalHandlers() {
	// Enable mouse inputs.
	core.App.TView.EnableMouse(true)

	// Set universal keybindings
	core.App.TView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlK: // Help page.
			ctrlKInput()
		case tcell.KeyCtrlS: // Search page.
			ctrlSInput()
		case tcell.KeyCtrlC: // Ctrl-C interrupt.
			ctrlCInput()
		}
		return event // Forward the event to the actual current primitive.
	})
}

// ctrlKInput : Shows the help page to the user.
func ctrlKInput() {
}

// ctrlSInput : Shows search page to the user.
// func ctrlSInput() {
// 	// Do not allow when on login screen.
// 	if page, _ := core.App.PageHolder.GetFrontPage(); page == utils.LoginPageID {
// 		return
// 	}
// 	// ShowSearchPage()
// }

// ctrlCInput : Sends an interrupt signal to the application to stop.
func ctrlCInput() {
	log.Println("TView stopped by Ctrl-C interrupt.")
	core.App.TView.Stop()
}

// ctrlSInput : Shows search page to the user.
func ctrlSInput() {
	// Do not allow when on login screen.
	// if page, _ := core.App.PageHolder.GetFrontPage(); page == utils.LoginPageID {
	// 	return
	// }
	// ShowSearchPage()
}

// setHandlers : Set handlers for the main page.
func (p *MainPage) setHandlers(cancel context.CancelFunc) {
	// Set table input captures.
	p.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		var reload bool
		switch event.Key() {
		// User wants to go to the next offset page.
		case tcell.KeyCtrlF:
			if p.CurrentOffset+limit >= p.MaxOffset {
				modal := okModal(utils.OffsetErrorModalID, "No more results to show.")
				ShowModal(utils.OffsetErrorModalID, modal)
			} else {
				// Update the new offset
				p.CurrentOffset += limit
			}
			reload = true
		case tcell.KeyCtrlB:
			if p.CurrentOffset == 0 {
				modal := okModal(utils.OffsetErrorModalID, "Already on first page.")
				ShowModal(utils.OffsetErrorModalID, modal)
			}
			reload = true
			// Update the new offset
			p.CurrentOffset = int(math.Max(0, float64(p.CurrentOffset-limit)))
		}

		if reload {
			// Cancel any current loading, and create a new one.
			cancel()
			go p.setLatestUpdatedAnimeTable()
		}
		return event
	})

	// Set table selected function.
	p.Table.SetSelectedFunc(func(row, _ int) {
		log.Printf("Selected row %d on main page.\n", row)
		animeRef := p.Table.GetCell(row, 0).GetReference()
		if animeRef == nil {
			return
		} else if anime, ok := animeRef.(*angoslayer.Anime); ok {
			ShowAnimePage(anime)
		}
	})
}
