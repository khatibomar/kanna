package service

import (
	"log"

	"github.com/khatibomar/kanna/app/core"
	"github.com/khatibomar/kanna/app/ui"
	"github.com/khatibomar/tohru"
	"github.com/rivo/tview"
)

func Start() {
	core.App = &core.Kanna{
		Client:     &tohru.TohruClient{},
		TView:      tview.NewApplication(),
		PageHolder: tview.NewPages(),
	}
	core.App.Initialise()
	cfg := tohru.NewConfig(core.App.Config.ClientID, core.App.Config.ClientSecret)
	core.App.Client = tohru.NewTohruClient(cfg)

	ui.ShowMainPage()
	log.Println("Initialised starting screen.")
	ui.SetUniversalHandlers()

	log.Println("Running app...")
	if err := core.App.TView.Run(); err != nil {
		log.Println(err)
		return
	}
}

func Shutdown() {
	core.App.Shutdown()
}
