package main

import "github.com/khatibomar/kanna/app/service"

func main() {
	service.Start()
	defer service.Shutdown()
}
