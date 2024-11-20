package database

import (
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/skutaada/gonergy/models"
)

func InsertEnergySpots(e *[]models.EnergySpot) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into energy(date, value) values (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, spot := range *e {
		_, err = stmt.Exec(spot.Date.Format(time.RFC3339), spot.Value)
		if err != nil {
			if errSql, ok := err.(sqlite3.Error); ok && errSql.ExtendedCode == sqlite3.ErrConstraintUnique {
				continue
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func SelectEnergySpotsForDay(date time.Time) ([]models.EnergySpot, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)
	rows, err := DB.Query("select * from energy where date between ? and ?", start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var energySpots []models.EnergySpot
	for rows.Next() {
		var spot models.EnergySpot
		var dateRaw string
		if err := rows.Scan(&spot.ID, &dateRaw, &spot.Value); err != nil {
			return nil, err
		}
		spot.Date, err = time.Parse(time.RFC3339, dateRaw)
		if err != nil {
			return nil, err
		}
		energySpots = append(energySpots, spot)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return energySpots, nil
}

func SelectEnergySpotCurrent(date time.Time) (*models.EnergySpot, error) {
	targetTime := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		date.Hour(),
		date.Minute() / 15 * 15,
		0,
		0,
		date.Location(),
	).Format(time.RFC3339)

	var energySpot models.EnergySpot
	var dateRaw string

	if err := DB.QueryRow("select * from energy where date = ?", targetTime).Scan(&energySpot.ID, &dateRaw, &energySpot.Value); err != nil {
		return nil, err
	}
	var err error
	energySpot.Date, err = time.Parse(time.RFC3339, dateRaw)
	if err != nil {
		return nil, err
	}

	return &energySpot, nil
}