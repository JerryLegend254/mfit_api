package main

import (
	"log"
)

func main() {

	app := &application{}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
