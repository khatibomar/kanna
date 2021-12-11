package main

import (
	"fmt"
	"log"
	"os"

	"github.com/khatibomar/tkanna/config"
)

func main() {
	f, err := os.Open(".ENV.toml")
	if err != nil {
		log.Fatalln(err)
	}
	c, err := config.New(f)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(c)
}
