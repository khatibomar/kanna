package core

import (
	"log"
	"os"

	"github.com/khatibomar/tohru"
	"github.com/rivo/tview"
)

var App *Kanna

type Kanna struct {
	Client *tohru.TohruClient

	TView      *tview.Application
	PageHolder *tview.Pages

	Config  *Config
	LogFile *os.File
}

func (t *Kanna) Initialise() error {
	if err := t.setUpLogging(); err != nil {
		log.Println("Unable to set up logging...")
		return err
	}

	if err := t.loadConfiguration(); err != nil {
		log.Println("Unable to read configuration file. Is it formatted correctly?")
		log.Println("If in doubt, delete the configuration file to start over!\n\nDetails:")
		return err
	}

	t.TView.SetRoot(t.PageHolder, true).SetFocus(t.PageHolder)

	return nil
}

func (t *Kanna) Shutdown() {
	App.TView.Sync()
	App.TView.Stop()

	if err := t.stopLogging(); err != nil {
		log.Println("Error while closing log file!")
	}
	log.Println("Application shutdown.")
}
