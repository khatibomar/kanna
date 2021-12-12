package main

import (
	"log"
	"os"

	"github.com/khatibomar/tkanna/config"
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
		errLogger.Fatalln(err)
	}
	infoLogger.Println("starting app...")
	infoLogger.Println(c)
}
