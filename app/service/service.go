package service

import (
	"log"

	"github.com/khatibomar/angoslayer"
	"github.com/khatibomar/tkanna/app/core"
	"github.com/khatibomar/tkanna/app/ui"
	"github.com/rivo/tview"
)

func Start() {
	core.App = &core.Tkanna{
		Client:     &angoslayer.AngoClient{},
		TView:      tview.NewApplication(),
		PageHolder: tview.NewPages(),
	}
	core.App.Initialise()
	cfg := angoslayer.NewConfig(core.App.Config.ClientID, core.App.Config.ClientSecret)
	core.App.Client = angoslayer.NewAngoClient(cfg)

	ui.ShowMainPage()
	log.Println("Initialised starting screen.")
	// ui.SetUniversalHandlers()

	log.Println("Running app...")
	if err := core.App.TView.Run(); err != nil {
		log.Fatalln(err)
	}
}

func Shutdown() {
	core.App.Shutdown()
}
