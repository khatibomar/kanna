package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/khatibomar/tkanna/config"
	"github.com/khatibomar/tkanna/internal/components/browser"
)

func main() {
	infoLogger := log.New(os.Stdout, "Info : ", log.Ldate|log.Ltime)
	errLogger := log.New(os.Stderr, "Error : ", log.Ldate|log.Ltime)
	f, err := os.Open(".ENV.toml")
	if err != nil {
		log.Fatalln(err)
	}
	c, err := config.New(f)
	if err != nil {
		log.Fatalln(err)
	}
	infoLogger.Println("starting app...")
	p := tea.NewProgram(browser.InitialModel(c))
	p.EnterAltScreen()
	if err := p.Start(); err != nil {
		errLogger.Fatalln(err)
	}
}
