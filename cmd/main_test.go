package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"pizza-hub/models"
	"strconv"
	"testing"
	"time"
)

func TestFlow(t *testing.T) {
	var resp models.Response
	//get menu
	req := httptest.NewRequest(http.MethodGet, "/menus", nil)
	w := httptest.NewRecorder()
	ListMenus(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("status code %s", res.Status)
		return
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	t.Logf("Menu: %v\n", string(data))
	json.Unmarshal(data, &resp)

	//order up
	var jsonStr = []byte(`{"customer":"arga","menu_id":[1,2]}`)
	reqOrder := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(jsonStr))
	wOrder := httptest.NewRecorder()
	AddOrder(wOrder, reqOrder)
	resOrder := wOrder.Result()
	if resOrder.StatusCode != http.StatusAccepted {
		t.Errorf("status code %s", resOrder.Status)
		return
	}
	defer resOrder.Body.Close()
	dataOrder, errOrder := io.ReadAll(resOrder.Body)
	if errOrder != nil {
		t.Errorf("expected error to be nil got %v", errOrder)
		return
	}
	var order models.Order
	t.Logf("Order: %v\n", string(dataOrder))
	json.Unmarshal(dataOrder, &resp)
	byteRes, _ := json.Marshal(resp.Data)
	json.Unmarshal(byteRes, &order)
	if order.ChefId == 0 {
		t.Log("chef is not assign yet")
	}
	//wait for chef to be assign
	time.Sleep(2 * time.Second)
	id := strconv.Itoa(order.Id)
	reqCheck := httptest.NewRequest(http.MethodGet, "/order/"+id, nil)
	wCheck := httptest.NewRecorder()
	GetOrderById(wCheck, reqCheck)
	resCheck := wCheck.Result()
	if resCheck.StatusCode != http.StatusOK {
		t.Errorf("status code %s", resCheck.Status)
		return
	}
	defer resCheck.Body.Close()
	dataCheck, errCheck := io.ReadAll(resCheck.Body)
	if errCheck != nil {
		t.Errorf("expected error to be nil got %v", errCheck)
		return
	}
	t.Logf("Order: %v\n", string(dataCheck))

	var orderNew models.Order
	json.Unmarshal(dataCheck, &resp)
	byteRes, _ = json.Marshal(resp.Data)
	json.Unmarshal(byteRes, &orderNew)

	if orderNew.ChefId != 0 {
		t.Logf("chef %v is assigned for order no %v", orderNew.ChefId, orderNew.Id)
	}

	//get chef status
	chefid := strconv.Itoa(orderNew.ChefId)
	reqChef := httptest.NewRequest(http.MethodGet, "/chef/"+chefid, nil)
	wChef := httptest.NewRecorder()
	GetChef(wChef, reqChef)
	resChef := wChef.Result()
	defer resChef.Body.Close()
	if resChef.StatusCode != http.StatusOK {
		t.Errorf("status code %s", resChef.Status)
		return
	}
	dataChef, errChef := io.ReadAll(resChef.Body)
	if errChef != nil {
		t.Errorf("expected error to be nil got %v", errChef)
		return

	}
	t.Logf("Chef: %v\n", string(dataChef))

	var chef models.Chef
	json.Unmarshal(dataChef, &resp)
	byteRes, _ = json.Marshal(resp.Data)
	json.Unmarshal(byteRes, &chef)
	if chef.IsOccupied {
		t.Log("chef is still busy")
		return
	} else {
		t.Log("chef is free")
	}
}
