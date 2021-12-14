package main

import "github.com/khatibomar/tkanna/app/service"

func main() {
	service.Start()
	defer service.Shutdown()
}
