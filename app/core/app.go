package core

import (
	"log"
	"os"

	"github.com/khatibomar/angoslayer"
	"github.com/rivo/tview"
)

var App *Tkanna

type Tkanna struct {
	Client *angoslayer.AngoClient

	TView      *tview.Application
	PageHolder *tview.Pages

	Config  *Config
	LogFile *os.File
}

func (t *Tkanna) Initialise() error {
	if err := t.setUpLogging(); err != nil {
		log.Fatalln("Unable to set up logging...")
	}

	if err := t.loadConfiguration(); err != nil {
		log.Println("Unable to read configuration file. Is it formatted correctly?")
		log.Println("If in doubt, delete the configuration file to start over!\n\nDetails:")
		log.Fatalln(err)
	}

	t.TView.SetRoot(t.PageHolder, true).SetFocus(t.PageHolder)

	return nil
}

func (t *Tkanna) Shutdown() {
	App.TView.Sync()
	App.TView.Stop()

	if err := t.stopLogging(); err != nil {
		log.Println("Error while closing log file!")
	}
	log.Println("Application shutdown.")
}
