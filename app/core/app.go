package core

import (
	"log"
	"os"

	"codeberg.org/omarkhatib/fafnir"
	"codeberg.org/omarkhatib/fafnir/repository"
	"codeberg.org/omarkhatib/tohru"
	"github.com/rivo/tview"
)

type Kanna struct {
	Client *tohru.TohruClient
	Fafnir *fafnir.Fafnir

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

	repo := repository.NewInMemoryRepo()

	fafnirCfg := fafnir.Config{
		ErrChan:                make(chan error, 100),
		Repo:                   repo,
		MaxConcurrentDownloads: t.Config.MaxConcurrentDownloads,
	}
	var err error
	t.Fafnir, err = fafnir.New(&fafnirCfg)

	if err != nil {
		return err
	}

	t.TView.SetRoot(t.PageHolder, true).SetFocus(t.PageHolder)
	return nil
}

func (t *Kanna) Shutdown() {
	t.TView.Sync()
	t.TView.Stop()

	if err := t.stopLogging(); err != nil {
		log.Println("Error while closing log file!")
	}
	log.Println("Application shutdown.")
}
