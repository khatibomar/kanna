package ui

import (
	"log"

	"github.com/khatibomar/kanna/app/core"
	"github.com/khatibomar/kanna/app/ui/utils"
	"github.com/rivo/tview"
)

// ShowModal : Make the app show a modal.
func ShowModal(id string, modal *tview.Modal) {
	core.App.TView.SetFocus(modal)
	core.App.PageHolder.AddPage(id, modal, true, true)
}

// okModal : Creates a new modal with an "OK" acknowledgement button.
func okModal(id, text string) *tview.Modal {
	modal := tview.NewModal()

	// Set modal attributes
	modal.SetText(text).
		SetBackgroundColor(utils.ModalColor).
		AddButtons([]string{"OK"}).
		SetFocus(0).
		SetDoneFunc(func(_ int, _ string) {
			core.App.PageHolder.RemovePage(id)
		})
	return modal
}

// confirmModal : Creates a new modal for confirmation.
// The user specifies the function to do when confirming.
// If the user cancels, then the modal is removed from the view.
func confirmModal(id, text, confirmButton string, f func()) *tview.Modal {
	// Create new modal
	modal := tview.NewModal()

	// Set modal attributes
	modal.SetText(text).
		SetBackgroundColor(utils.ModalColor).
		AddButtons([]string{confirmButton, "Cancel"}).
		SetFocus(0).
		SetDoneFunc(func(buttonIndex int, _ string) {
			if buttonIndex == 0 {
				f()
			}
			log.Printf("Removing %s modal\n", id)
			core.App.PageHolder.RemovePage(id)
		})
	return modal
}

// confirmModal : Creates a new modal for confirmation.
// The user specifies the function to do when confirming.
// If the user cancels, then the modal is removed from the view.
func watchOrDownloadModal(id, text string, stream func(chan error, chan string), download func(chan error, chan string), errChan chan error, infoChan chan string) *tview.Modal {
	// Create new modal
	modal := tview.NewModal()

	// Set modal attributes
	modal.SetText(text).
		SetBackgroundColor(utils.ModalColor).
		AddButtons([]string{"Stream", "Download", "Cancel"}).
		SetFocus(0).
		SetDoneFunc(func(buttonIndex int, _ string) {
			if buttonIndex == 0 {
				stream(errChan, infoChan)
			} else if buttonIndex == 1 {
				download(errChan, infoChan)
			}
			log.Printf("Removing %s modal\n", id)
			core.App.PageHolder.RemovePage(id)
		})
	return modal
}
