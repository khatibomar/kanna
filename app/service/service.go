package service

import (
	"log"

	"github.com/khatibomar/angoslayer"
	"github.com/khatibomar/tkanna/app/core"
	"github.com/rivo/tview"
)

func Start() {
	cfg := angoslayer.NewConfig(core.App.Config.ClientID, core.App.Config.ClientSecret)
	core.App = &core.Tkanna{
		Client:     angoslayer.NewAngoClient(cfg),
		TView:      tview.NewApplication(),
		PageHolder: tview.NewPages(),
	}

	// ui.ShowMainPage()
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
