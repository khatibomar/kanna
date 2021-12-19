package service

import (
	"log"

	"github.com/khatibomar/kanna/app/core"
	"github.com/khatibomar/kanna/app/ui"
	"github.com/khatibomar/tohru"
	"github.com/rivo/tview"
)

func Start() {
	core := &core.Kanna{
		Client:     &tohru.TohruClient{},
		TView:      tview.NewApplication(),
		PageHolder: tview.NewPages(),
	}
	core.Initialise()
	cfg := tohru.NewConfig(core.Config.ClientID, core.Config.ClientSecret)
	core.Client = tohru.NewTohruClient(cfg)

	ui.ShowMainPage(core)
	log.Println("Initialised starting screen.")
	ui.SetUniversalHandlers(core)

	log.Println("Running app...")
	if err := core.TView.Run(); err != nil {
		log.Println(err)
		return
	}
	defer core.Shutdown()
}
