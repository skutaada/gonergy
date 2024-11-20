package models

import "time"

type EnergySpot struct {
	ID int
	Date time.Time `json:"date"`
	Value float32 `json:"value"`
}

type EnergyDay struct {
	Data []EnergySpot `json:"data"`
}