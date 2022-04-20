package ui

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"strconv"

	"codeberg.org/omarkhatib/kanna/app/core"
	"codeberg.org/omarkhatib/kanna/app/ui/utils"
	"codeberg.org/omarkhatib/tohru"
	"github.com/gdamore/tcell/v2"
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

		dwnF := func(selected map[int]struct{}, errChan chan error, infoChan chan string) {
			var selection EpisodeSelection
			var ok bool
			for index := range selected {
				if selection, ok = p.Table.GetCell(index, 0).GetReference().(EpisodeSelection); !ok {
					return
				}
				log.Printf("Downloading episode %s\n", selection.episode.EpisodeName)
				animeID, _ := strconv.Atoi(selection.animeID)
				episodeID, _ := strconv.Atoi(selection.episode.EpisodeID)
				episode, err := p.Core.Client.EpisodeService.GetEpisodeDetails(animeID, episodeID)
				if err != nil {
					log.Printf("Failed to get episode: %s\n", err.Error())
					return
				}
				go p.saveEpisode(&episode, errChan, infoChan)
			}
			info := fmt.Sprintf("Download Starting... \nyou can find file in %s\nif error happened it will be reported", p.Core.Config.DownloadDir)
			modal := okModal(p.Core, utils.InfoModalID, info)
			ShowModal(p.Core, utils.InfoModalID, modal)
		}

		streamF := func(selected map[int]struct{}, errChan chan error, infoChan chan string) {
			var selection EpisodeSelection
			var ok bool
			for index := range selected {
				if selection, ok = p.Table.GetCell(index, 0).GetReference().(EpisodeSelection); !ok {
					return
				}
				log.Printf("Streaming episode %s\n", selection.episode.EpisodeName)
				animeID, _ := strconv.Atoi(selection.animeID)
				episodeID, _ := strconv.Atoi(selection.episode.EpisodeID)
				episode, err := p.Core.Client.EpisodeService.GetEpisodeDetails(animeID, episodeID)
				if err != nil {
					log.Printf("Failed to get episode: %s\n", err.Error())
					return
				}
				go p.streamEpisode(&episode, errChan)
			}

			log.Println(selected)
			modal := okModal(p.Core, utils.InfoModalID, "Stream Starting...\n this operation may take few minutes based on internet connection and mpv launch \nif error happened it will be reported")
			ShowModal(p.Core, utils.InfoModalID, modal)
		}

		selected := p.sWrap.CopySelection()
		if len(selected) > 1 {
			modal := confirmDownloadModal(p.Core, selected, dwnF, errChan, infoChan)
			ShowModal(p.Core, utils.DownloadModalID, modal)
		} else {
			modal := watchOrDownloadModal(p.Core, utils.WatchOrDownloadModalID, "Select Option", selected, streamF, dwnF, errChan, infoChan)
			ShowModal(p.Core, utils.WatchOrDownloadModalID, modal)
		}

		// using hashing here is because if many errors or info reported
		// at same time closing all of them because they have same ID
		// so adding an hash will prevent that and make each okMadal
		// have it's own unique ID
		// TODO(khatibomar): find a lighter way to make hashes
		go func(errChan chan error, infoChan chan string) {
			for {
				select {
				case err := <-errChan:
					log.Println(err)
					p.Core.TView.QueueUpdateDraw(func() {
						hash := GetMD5Hash(err.Error())
						modal := okModal(p.Core, utils.GenericAPIErrorModalID+hash, err.Error())
						ShowModal(p.Core, utils.GenericAPIErrorModalID+hash, modal)
					})

				case info := <-infoChan:
					log.Println(info)
					p.Core.TView.QueueUpdateDraw(func() {
						hash := GetMD5Hash(info)
						modal := okModal(p.Core, utils.InfoModalID+hash, info)
						ShowModal(p.Core, utils.InfoModalID+hash, modal)
					})
				}
			}
		}(errChan, infoChan)

		for index := range selected {
			p.sWrap.RemoveSelection(index)
			p.markUnselected(index)
		}
	})

	// Set table input captures.
	p.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlE: // User selects this manga row.
			p.ctrlEInput()
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

func (p *AnimePage) ctrlEInput() {
	row, _ := p.Table.GetSelection()
	// If the row is already in the Selection, we deselect. Else, we add.
	if p.sWrap.HasSelection(row) {
		p.markUnselected(row)
	} else {
		p.markSelected(row)
	}
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

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
