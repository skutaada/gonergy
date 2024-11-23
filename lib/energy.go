package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/skutaada/gonergy/database"
	"github.com/skutaada/gonergy/models"
)

func FetchLatestEnergy() ([]models.EnergySpot, error) {
	res, err := http.Get("https://apis.smartenergy.at/market/v1/price")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var energyDay models.EnergyDay
	if err := json.Unmarshal(data, &energyDay); err != nil {
		return nil, err
	}

	return energyDay.Data, nil
}

func DailyFetchInsert() {
	for {
		now := time.Now()
		target := time.Date(now.Year(), now.Month(), now.Day(), 17, 15, 0, 0, now.Location())

		if now.After(target) {
			target = target.Add(24 * time.Hour)
		}
		waitDuration := target.Sub(now)
		timer := time.NewTimer(waitDuration)

		<-timer.C
		res, err := FetchLatestEnergy()
		if err != nil {
			fmt.Println(err)
		}
		if err := database.InsertEnergySpots(&res); err != nil {
			fmt.Println(err)
		}
		fmt.Println("Sucessfully saved the things")
	}
}

func DailyFetchInsert2(c <-chan bool, done chan<- bool) {
	for {
		now := time.Now()
		target := time.Date(now.Year(), now.Month(), now.Day(), 17, 15, 0, 0, now.Location())

		if now.After(target) {
			target = target.Add(24 * time.Hour)
		}
		waitDuration := target.Sub(now)
		timer := time.NewTimer(waitDuration)

		select {
		case <-timer.C:
			res, err := FetchLatestEnergy()
			if err != nil {
				fmt.Println(err)
			}
			if err := database.InsertEnergySpots(&res); err != nil {
				fmt.Println(err)
			}
			fmt.Println("Sucessfully saved the data from smartEnergy")
		case <-c:
			timer.Stop()
			done <- true
			fmt.Println("Gracefully shutdown Fetcher.")
			return
		}
	}
}
