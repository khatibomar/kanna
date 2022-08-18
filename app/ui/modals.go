package ui

import (
	"log"

	"codeberg.org/omarkhatib/kanna/app/core"
	"codeberg.org/omarkhatib/kanna/app/ui/utils"
	"github.com/rivo/tview"
)

type actionFunc func(map[int]struct{}, chan error, chan string)

// ShowModal : Make the app show a modal.
func ShowModal(core *core.Kanna, id string, modal *tview.Modal) {
	core.TView.SetFocus(modal)
	core.PageHolder.AddPage(id, modal, true, true)
}

// okModal : Creates a new modal with an "OK" acknowledgement button.
func okModal(core *core.Kanna, id, text string) *tview.Modal {
	modal := tview.NewModal()

	// Set modal attributes
	modal.SetText(text).
		SetBackgroundColor(utils.ModalColor).
		AddButtons([]string{"OK"}).
		SetFocus(0).
		SetDoneFunc(func(_ int, _ string) {
			core.PageHolder.RemovePage(id)
		})
	return modal
}

func confirmDownloadModal(core *core.Kanna, selected map[int]struct{}, f actionFunc, errChan chan error, infoChan chan string) *tview.Modal {
	// Create new modal
	modal := tview.NewModal()

	// Set modal attributes
	modal.SetText("Confirm download").
		SetBackgroundColor(utils.ModalColor).
		AddButtons([]string{"Download", "Cancel"}).
		SetFocus(0).
		SetDoneFunc(func(buttonIndex int, _ string) {
			if buttonIndex == 0 {
				f(selected, errChan, infoChan)
			}
			log.Printf("Removing %s modal\n", utils.DownloadModalID)
			core.PageHolder.RemovePage(utils.DownloadModalID)
		})
	return modal
}

// confirmModal : Creates a new modal for confirmation.
// The user specifies the function to do when confirming.
// If the user cancels, then the modal is removed from the view.

func watchOrDownloadModal(core *core.Kanna, id, text string, selected map[int]struct{}, stream actionFunc, download actionFunc, errChan chan error, infoChan chan string) *tview.Modal {
	// Create new modal
	modal := tview.NewModal()

	// Set modal attributes
	modal.SetText(text).
		SetBackgroundColor(utils.ModalColor).
		AddButtons([]string{"Stream", "Download", "Cancel"}).
		SetFocus(0).
		SetDoneFunc(func(buttonIndex int, _ string) {
			if buttonIndex == 0 {
				stream(selected, errChan, infoChan)
			} else if buttonIndex == 1 {
				download(selected, errChan, infoChan)
			}
			log.Printf("Removing %s modal\n", id)
			core.PageHolder.RemovePage(id)
		})
	return modal
}
