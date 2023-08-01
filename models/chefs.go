package models

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Chef struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	IsOccupied bool   `json:"is_occupied"`
}

func (c *Chef) Validate() error {
	if c.Name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

func (c *Chef) Bind(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(&c)
}

var ChefDB = []Chef{
	{
		Id:         1,
		Name:       "Ersa",
		IsOccupied: false,
	},
}

func (c *Chef) EditChef() {
	for _, v := range ChefDB {
		if v.Id == c.Id {
			ChefDB[c.Id-1] = *c
		}
	}
}
