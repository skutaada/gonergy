package main

import (
	"log"

	"github.com/skutaada/gonergy/database"
	"github.com/skutaada/gonergy/lib"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	go lib.DailyFetchInsert()
}
