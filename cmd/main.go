package main

import (
	"net/http"
	"pizza-hub/helpers"
	"pizza-hub/models"
	"strconv"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/chef", AddChef)
	http.HandleFunc("/chef/", GetChef)
	http.HandleFunc("/menus", ListMenus)
	http.HandleFunc("/orders", AddOrder)
	http.HandleFunc("/order/", GetOrderById)
	http.HandleFunc("/all/order", ListOrders)

	http.ListenAndServe(":3000", nil)
}

func AddChef(w http.ResponseWriter, r *http.Request) {
	ok := helpers.MethodHelper(r, "POST")
	if !ok {
		models.HandleResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}
	var newChef models.Chef
	if err := newChef.Bind(r); err != nil {
		models.HandleResponse(w, http.StatusInternalServerError, "cannot get body", nil)
		return
	}

	if err := newChef.Validate(); err != nil {
		models.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	newChef.Id = models.ChefDB[len(models.ChefDB)-1].Id + 1
	newChef.IsOccupied = false

	models.ChefDB = append(models.ChefDB, newChef)

	models.HandleResponse(w, http.StatusAccepted, "new chef created", newChef)
}

func GetChef(w http.ResponseWriter, r *http.Request) {
	ok := helpers.MethodHelper(r, "GET")
	if !ok {
		models.HandleResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/chef/")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		models.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	var resChef *models.Chef
	for _, v := range models.ChefDB {
		if v.Id == idInt {
			resChef = &models.Chef{
				Id:         idInt,
				Name:       v.Name,
				IsOccupied: v.IsOccupied,
			}
		}
	}
	if resChef == nil {
		models.HandleResponse(w, http.StatusNotFound, "chef not found", nil)
		return
	}
	models.HandleResponse(w, http.StatusOK, "chef found", resChef)
}

func ListMenus(w http.ResponseWriter, r *http.Request) {
	ok := helpers.MethodHelper(r, "GET")
	if !ok {
		models.HandleResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}
	if models.PizzaMenuDB == nil {
		models.HandleResponse(w, http.StatusNotFound, "menu not found", nil)
		return
	}
	var response []models.MenuResp
	for _, v := range models.PizzaMenuDB {
		response = append(response, models.MenuResp{
			Id:             v.Id,
			Name:           v.Name,
			DurationString: v.DurationString,
		})
	}
	models.HandleResponse(w, http.StatusOK, "get list menus", response)
}

func AddOrder(w http.ResponseWriter, r *http.Request) {
	ok := helpers.MethodHelper(r, "POST")
	if !ok {
		models.HandleResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}
	var newOrder models.Order
	if err := newOrder.Bind(r); err != nil {
		models.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if err := newOrder.Validate(); err != nil {
		models.HandleResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var eta time.Duration
	for i, v := range newOrder.MenuId {
		if i <= len(models.PizzaMenuDB) {
			for _, x := range models.PizzaMenuDB {
				if x.Id == v {
					eta = eta + x.Duration
				}
			}
		} else {
			models.HandleResponse(w, http.StatusNotFound, "menu not found", nil)
			return
		}
	}
	if len(models.OrderDB) == 0 {
		newOrder.Id = 1
	} else {
		newOrder.Id = models.OrderDB[len(models.OrderDB)-1].Id + 1
	}

	models.OrderDB = append(models.OrderDB, newOrder)

	go newOrder.ProcessOrder(eta)

	models.HandleResponse(w, http.StatusAccepted, "order created", newOrder)
}

func GetOrderById(w http.ResponseWriter, r *http.Request) {
	ok := helpers.MethodHelper(r, "GET")
	if !ok {
		models.HandleResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/order/")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		models.HandleResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	var resOrder *models.Order
	for _, v := range models.OrderDB {
		if v.Id == idInt {
			resOrder = &models.Order{
				Id:            v.Id,
				Customer:      v.Customer,
				MenuId:        v.MenuId,
				ChefId:        v.ChefId,
				EstimatedTime: v.EstimatedTime,
				Done:          v.Done,
			}
		}
	}
	if resOrder == nil {
		models.HandleResponse(w, http.StatusNotFound, "order not found", nil)
		return
	}
	models.HandleResponse(w, http.StatusOK, "order found", resOrder)
}

func ListOrders(w http.ResponseWriter, r *http.Request) {
	ok := helpers.MethodHelper(r, "GET")
	if !ok {
		models.HandleResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}
	if models.OrderDB == nil {
		models.HandleResponse(w, http.StatusNotFound, "order not found", nil)
	}
	var response []models.Order
	for _, v := range models.OrderDB {
		response = append(response, models.Order{
			Id:            v.Id,
			Customer:      v.Customer,
			MenuId:        v.MenuId,
			ChefId:        v.ChefId,
			EstimatedTime: v.EstimatedTime,
			Done:          v.Done,
		})
	}
	models.HandleResponse(w, http.StatusOK, "get list order", response)
}
