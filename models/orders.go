package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Order struct {
	Id            int    `json:"id"`
	Customer      string `json:"customer"`
	MenuId        []int  `json:"menu_id"`
	ChefId        int    `json:"chef_id"`
	EstimatedTime string `json:"estimated_time,omitempty"`
	Done          bool   `json:"done"`
}

func (o *Order) Validate() error {
	if o.Customer == "" {
		return errors.New("customer cannot be empty")
	}
	if o.MenuId == nil {
		return errors.New("menu cannot be empty")
	}
	return nil
}

func (o *Order) Bind(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(&o)
}

var OrderDB = []Order{}

func (o *Order) ProcessOrder(eta time.Duration) {
	var editChef Chef
	for o.ChefId == 0 {
		o.AssignChef()
		if o.ChefId != 0 {
			o.EstimatedTime = time.Now().Add(eta).Format("2006-01-02 15:04:05")
			o.EditOrder()
			editChef.Id = o.ChefId
			editChef.IsOccupied = true
			editChef.EditChef()
		}
	}
	time.Sleep(eta)
	o.Done = true
	o.EditOrder()
	editChef.IsOccupied = false
	editChef.EditChef()
	fmt.Printf("order no %v is finished", o.Id)
}

func (o *Order) EditOrder() {
	for _, v := range OrderDB {
		if v.Id == o.Id {
			OrderDB[o.Id-1] = *o
		}
	}
}

func (o *Order) AssignChef() {
	for _, v := range ChefDB {
		if !v.IsOccupied {
			o.ChefId = v.Id
		}
	}
}
