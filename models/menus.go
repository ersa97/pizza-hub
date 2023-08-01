package models

import "time"

type Menu struct {
	Id             int           `json:"id"`
	Name           string        `json:"name"`
	Duration       time.Duration `json:"duration"`
	DurationString string        `json:"duration_string"`
}
type MenuResp struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	DurationString string `json:"duration_string"`
}

var PizzaMenuDB = []Menu{
	{
		Id:             1,
		Name:           "Pizza Cheese",
		Duration:       time.Duration(3 * time.Minute),
		DurationString: "3 minutes",
	},
	{
		Id:             2,
		Name:           "Pizza BBQ",
		Duration:       time.Duration(5 * time.Minute),
		DurationString: "5 minutes",
	},
}
