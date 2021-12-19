package ui

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/khatibomar/kanna/app/core"
	"github.com/khatibomar/kanna/app/ui/utils"
	"github.com/khatibomar/tohru"
)

// SetUniversalHandlers : Set universal inputs for the app.
func SetUniversalHandlers(core *core.Kanna) {
	// Enable mouse inputs.
	core.TView.EnableMouse(true)

	// Set universal keybindings
	core.TView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlK: // Help page.
			ctrlKInput(core)
		case tcell.KeyCtrlS: // Search page.
			ctrlSInput(core)
		case tcell.KeyCtrlC: // Ctrl-C interrupt.
			ctrlCInput(core)
		}
		return event // Forward the event to the actual current primitive.
	})
}

// ctrlKInput : Shows the help page to the user.
func ctrlKInput(core *core.Kanna) {
	ShowHelpPage(core)
}

// ctrlSInput : Shows search page to the user.
func ctrlSInput(core *core.Kanna) {
	ShowSearchPage(core)
}

// ctrlCInput : Sends an interrupt signal to the application to stop.
func ctrlCInput(core *core.Kanna) {
	log.Println("TView stopped by Ctrl-C interrupt.")
	core.TView.Stop()
}

// setHandlers : Set handlers for the main page.
func (p *MainPage) setHandlers(cancel context.CancelFunc, searchParams *SearchParams) {
	// Set table input captures.
	p.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		var reload bool
		switch event.Key() {
		// User wants to go to the next offset page.
		case tcell.KeyCtrlF:
			if p.CurrentOffset+limit >= maxOffset {
				modal := okModal(p.Core, utils.OffsetErrorModalID, "No more results to show.")
				ShowModal(p.Core, utils.OffsetErrorModalID, modal)
			} else {
				// Update the new offset
				p.CurrentOffset += limit
			}
			reload = true
		case tcell.KeyCtrlB:
			if p.CurrentOffset == 0 {
				modal := okModal(p.Core, utils.OffsetErrorModalID, "Already on first page.")
				ShowModal(p.Core, utils.OffsetErrorModalID, modal)
			}
			reload = true
			// Update the new offset
			p.CurrentOffset = int(math.Max(0, float64(p.CurrentOffset-limit)))
		}

		if reload {
			// Cancel any current loading, and create a new one.
			cancel()
			if searchParams != nil {
				go p.setLatestUpdatedAnimeTable(searchParams)
			} else {
				go p.setLatestUpdatedAnimeTable(nil)
			}
		}
		return event
	})

	// Set table selected function.
	p.Table.SetSelectedFunc(func(row, _ int) {
		log.Printf("Selected row %d on main page.\n", row)
		animeRef := p.Table.GetCell(row, 0).GetReference()
		if animeRef == nil {
			return
		} else if anime, ok := animeRef.(*tohru.Anime); ok {
			ShowAnimePage(p.Core, anime)
		}
	})
}

// setHandlers : Set handlers for the help page.
func (p *HelpPage) setHandlers() {
	// Set grid input captures.
	p.Grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			p.Core.PageHolder.RemovePage(utils.HelpPageID)
		}
		return event
	})
}

// setHandlers : Set handlers for the page.
func (p *AnimePage) setHandlers(cancel context.CancelFunc) {
	// Set grid input captures.
	p.Grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			cancel()
			p.Core.PageHolder.RemovePage(utils.AnimePageID)
		}
		return event
	})

	// Set table selected function.
	p.Table.SetSelectedFunc(func(row, _ int) {
		// We add the current Selection if the there are no selected rows currently.
		if !p.sWrap.HasSelections() {
			p.sWrap.AddSelection(row)
		}
		log.Println("Creating and showing confirm download modal...")

		errChan := make(chan error)
		infoChan := make(chan string)

		selected := p.sWrap.CopySelection()
		dwnF := func(errChan chan error, infoChan chan string) {
			var episode *tohru.Episode
			var ok bool
			for index := range selected {
				if episode, ok = p.Table.GetCell(index, 0).GetReference().(*tohru.Episode); !ok {
					return
				}
				go p.saveEpisode(episode, errChan, infoChan)
			}
			for index := range selected {
				p.sWrap.RemoveSelection(index)
			}
			info := fmt.Sprintf("Download Starting... \nyou can find file in %s\nif error happened it will be reported", p.Core.Config.DownloadDir)
			modal := okModal(p.Core, utils.InfoModalID, info)
			ShowModal(p.Core, utils.InfoModalID, modal)
		}

		streamF := func(errChan chan error, infoChan chan string) {
			var episode *tohru.Episode
			var ok bool
			for index := range selected {
				if episode, ok = p.Table.GetCell(index, 0).GetReference().(*tohru.Episode); !ok {
					return
				}
				log.Printf("Streaming episode %s\n", episode.EpisodeName)
				log.Println(episode.EpisodeUrls)
				go p.streamEpisode(episode, errChan)
			}

			log.Println(selected)
			for index := range selected {
				p.sWrap.RemoveSelection(index)
			}
			modal := okModal(p.Core, utils.InfoModalID, "Stream Starting...\n this operation may take few minutes based on internet connection and mpv launch \nif error happened it will be reported")
			ShowModal(p.Core, utils.InfoModalID, modal)
		}

		if len(p.sWrap.Selection) > 1 {
			modal := confirmModal(p.Core, utils.DownloadModalID, "Download episode(s)?", "Yes", dwnF, errChan, infoChan)
			ShowModal(p.Core, utils.DownloadModalID, modal)
		} else {
			log.Println(selected)
			for index := range selected {
				p.sWrap.RemoveSelection(index)
			}
			modal := watchOrDownloadModal(p.Core, utils.WatchOrDownloadModalID, "Select Option", streamF, dwnF, errChan, infoChan)
			ShowModal(p.Core, utils.WatchOrDownloadModalID, modal)
		}
		go func(errChan chan error, infoChan chan string) {
			for {
				select {
				case err := <-errChan:
					log.Println(err)
					p.Core.TView.QueueUpdateDraw(func() {
						modal := okModal(p.Core, utils.GenericAPIErrorModalID, err.Error())
						ShowModal(p.Core, utils.GenericAPIErrorModalID, modal)
					})

				case info := <-infoChan:
					log.Println(info)
					p.Core.TView.QueueUpdateDraw(func() {
						modal := okModal(p.Core, utils.InfoModalID, info)
						ShowModal(p.Core, utils.InfoModalID, modal)
					})
				}
			}
		}(errChan, infoChan)
	})

	// Set table input captures.
	p.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		// case tcell.KeyCtrlE: // User selects this manga row.
		// 	p.ctrlEInput()
		case tcell.KeyCtrlA: // User wants to toggle select All.
			p.ctrlAInput()
			// case tcell.KeyCtrlR: // User wants to toggle read status for Selection.
			// 	p.ctrlRInput()
			// case tcell.KeyCtrlQ:
			// 	p.ctrlQInput()
		}
		return event
	})
}

func (p *AnimePage) ctrlAInput() {
	// Toggle Selection.
	p.markAll()
}

// setHandlers : Set handlers for the search page.
func (p *SearchPage) setHandlers() {
	// Set grid input captures.
	p.Grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc: // When user presses ESC, then we remove the Search page.
			p.Core.PageHolder.RemovePage(utils.SearchPageID)
		case tcell.KeyTab: // When user presses Tab, they are sent back to the search form.
			p.Core.TView.SetFocus(p.Form)
		}
		return event
	})

	// Set up input capture for the search bar.
	p.Form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown: // When user presses KeyDown, they are sent to the search results table.
			p.Core.TView.SetFocus(p.Table)
		}
		return event
	})
}
